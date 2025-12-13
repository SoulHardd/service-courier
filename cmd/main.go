package main

import (
	"avito/internal/config"
	courierHandler "avito/internal/handlers/courier"
	deliveryHandler "avito/internal/handlers/delivery"
	courierRepo "avito/internal/repository/courier"
	deliveryRepo "avito/internal/repository/delivery"
	"avito/internal/server"
	"avito/internal/useCase/courier"
	"avito/internal/useCase/delivery"
	"avito/pkg/connections"
	"context"
	"log"
	"os/signal"
	"syscall"
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

	srv := server.New(cfg.ServerCfg, pool, courierController, deliveryController)

	deliveryUseCase.StartMonitorDeliveries(appCtx, cfg.TimeCfg.DeliveryMonitorInterval)

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
