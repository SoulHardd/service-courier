package main

import (
	"avito/internal/config"
	"avito/internal/gateway/order"
	courierHandler "avito/internal/handlers/http/courier"
	deliveryHandler "avito/internal/handlers/http/delivery"
	"avito/internal/logger"
	courierRepo "avito/internal/repository/courier"
	deliveryRepo "avito/internal/repository/delivery"
	"avito/internal/server"
	"avito/internal/useCase/courier"
	"avito/internal/useCase/delivery"
	"avito/internal/worker"
	"avito/pkg/connections"
	pb "avito/proto/order"
	"context"
	"log"
	"os/signal"
	"syscall"

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
	defer zapLogger.Sync()
	lg := logger.NewZapLogger(zapLogger)
	srv := server.New(cfg.ServerCfg, pool, courierController, deliveryController, lg)

	deliveryUseCase.StartMonitorDeliveries(appCtx, cfg.TimeCfg.DeliveryMonitorInterval)

	conn, err := grpc.NewClient(cfg.GrpcCfg.OrderServiceGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	o := pb.NewOrdersServiceClient(conn)
	orderGateway := order.NewGateway(o)
	orderWorker := worker.New(orderGateway, deliveryUseCase, cfg.TimeCfg.OrderWorkerInterval)
	orderWorker.Run(appCtx)

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
