package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/SoulHardd/service-courier/internal/config"
	"github.com/SoulHardd/service-courier/internal/gateway/order"
	courierHandler "github.com/SoulHardd/service-courier/internal/handlers/http/courier"
	deliveryHandler "github.com/SoulHardd/service-courier/internal/handlers/http/delivery"
	"github.com/SoulHardd/service-courier/internal/logger"
	"github.com/SoulHardd/service-courier/internal/pprof"
	courierRepo "github.com/SoulHardd/service-courier/internal/repository/courier"
	deliveryRepo "github.com/SoulHardd/service-courier/internal/repository/delivery"
	"github.com/SoulHardd/service-courier/internal/server"
	"github.com/SoulHardd/service-courier/internal/useCase/courier"
	"github.com/SoulHardd/service-courier/internal/useCase/delivery"
	"github.com/SoulHardd/service-courier/internal/worker"
	"github.com/SoulHardd/service-courier/pkg/connections"
	pb "github.com/SoulHardd/service-courier/proto/order"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	pool := connections.InitPool(appCtx, cfg.DatabaseCfg)

	courierRepository := courierRepo.New(pool)
	courierUseCase := courier.New(courierRepository)
	courierController := courierHandler.New(courierUseCase)

	deliveryRepository := deliveryRepo.New(pool)
	deliveryUseCase := delivery.New(deliveryRepository, courierRepository)
	deliveryController := deliveryHandler.New(deliveryUseCase)

	zapLogger, _ := zap.NewProduction()
	defer func() {
		_ = zapLogger.Sync()
	}()
	lg := logger.NewZapLogger(zapLogger)
	srv := server.New(cfg.ServerCfg, pool, courierController, deliveryController, lg, cfg.RateLimiterCfg)

	deliveryUseCase.StartMonitorDeliveries(appCtx, cfg.TimeCfg.DeliveryMonitorInterval)

	conn, err := grpc.NewClient(cfg.GrpcCfg.OrderServiceGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()
	o := pb.NewOrdersServiceClient(conn)
	orderGateway := order.NewGateway(o, cfg.GrpcCfg)
	orderWorker := worker.New(orderGateway, deliveryUseCase, cfg.TimeCfg.OrderWorkerInterval)
	orderWorker.Run(appCtx)

	pprofServer := pprof.New(":6060")
	pprofServer.Run()

	srv.ListenAndServe()

	<-shutdownCtx.Done()
	log.Printf("Shutting down service-courier")
	gsCtx, gsCancel := context.WithTimeout(context.Background(), cfg.TimeCfg.ShutdownTimeout)
	defer gsCancel()

	if err := srv.HttpSrv.Shutdown(gsCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	pool.Close()
	log.Println("Graceful shutdown complete")
}
