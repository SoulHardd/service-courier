package main

import (
	"avito/internal/config"
	courierHandler "avito/internal/handlers/courier"
	courierRepo "avito/internal/repository/courier"
	"avito/internal/server"
	"avito/internal/useCase/courier"
	"avito/pkg/connections"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
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

	srv := server.New(cfg.ServerCfg, pool, courierController)

	srv.ListenAndServe()

	<-shutdownCtx.Done()
	log.Printf("Shutting down service-courier")
	gsCtx, gsCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer gsCancel()

	if err := srv.HttpSrv.Shutdown(gsCtx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}

	pool.Close()
	log.Println("Graceful shutdown complete")
}
