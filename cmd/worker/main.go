package main

import (
	"avito/internal/config"
	"avito/internal/gateway/order"
	changedHandler "avito/internal/handlers/queue/order/changed"
	courierRepo "avito/internal/repository/courier"
	deliveryRepo "avito/internal/repository/delivery"
	"avito/internal/useCase/courier"
	"avito/internal/useCase/delivery"
	changedUc "avito/internal/useCase/order/changed"
	"avito/pkg/connections"
	pb "avito/proto/order"
	"context"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

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
	defer pool.Close()

	courierRepository := courierRepo.New(pool)
	courierUseCase := courier.New(courierRepository)
	deliveryRepository := deliveryRepo.New(pool)
	deliveryUseCase := delivery.New(deliveryRepository, courierRepository)

	conn, err := grpc.NewClient(cfg.GrpcCfg.OrderServiceGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()
	o := pb.NewOrdersServiceClient(conn)
	orderGateway := order.NewGateway(o)

	orderUseCase := changedUc.NewOrderUseCase(orderGateway, deliveryUseCase, courierUseCase)

	handler := changedHandler.New(orderUseCase)

	saramaCfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(cfg.KafkaCfg.Version)
	if err != nil {
		log.Fatalf("invalid kafka version: %v", err)
	}
	saramaCfg.Version = version

	switch cfg.KafkaCfg.InitialOffset {
	case "oldest":
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		log.Fatalf("invalid initial offset: %s", cfg.KafkaCfg.InitialOffset)
	}

	saramaCfg.Consumer.Offsets.AutoCommit.Enable = cfg.KafkaCfg.AutoCommitEnable
	saramaCfg.Consumer.Offsets.AutoCommit.Interval = cfg.KafkaCfg.AutoCommitInterval

	kafkaClient, err := sarama.NewConsumerGroup(cfg.KafkaCfg.Brokers, cfg.KafkaCfg.GroupId, saramaCfg)
	if err != nil {
		log.Fatalf("Unabled to create kafka consumer group: %v", err)
	}
	defer kafkaClient.Close()

	go func() {
		for {
			if err := kafkaClient.Consume(appCtx, []string{cfg.KafkaCfg.OrderTopic}, handler); err != nil {
				log.Printf("kafka consume error: %v", err)
			}

			if appCtx.Err() != nil {
				return
			}
		}
	}()

	<-shutdownCtx.Done()
	log.Printf("Shutting down kafka worker")

	_, gsCancel := context.WithTimeout(context.Background(), cfg.TimeCfg.ShutdownTimeout)
	defer gsCancel()

	log.Println("Graceful shutdown kafka complete")
}
