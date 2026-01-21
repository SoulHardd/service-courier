package main

import (
	"avito/internal/config"
	"avito/internal/gateway/order"
	changedHandler "avito/internal/handlers/queue/order/changed"
	"avito/internal/logger"
	courierRepo "avito/internal/repository/courier"
	deliveryRepo "avito/internal/repository/delivery"
	"avito/internal/useCase/courier"
	"avito/internal/useCase/delivery"
	changedUc "avito/internal/useCase/order/changed"
	"avito/pkg/connections"
	pb "avito/proto/order"
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"os/signal"
	"syscall"
)

func main() {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	lg := logger.NewZapLogger(zapLogger)

	pool := connections.InitPool(appCtx, cfg.DatabaseCfg)
	defer pool.Close()

	courierRepository := courierRepo.New(pool)
	courierUseCase := courier.New(courierRepository)
	deliveryRepository := deliveryRepo.New(pool)
	deliveryUseCase := delivery.New(deliveryRepository, courierRepository)

	conn, err := grpc.NewClient(cfg.GrpcCfg.OrderServiceGrpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		lg.Fatal("Failed to connect to gRPC server",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "address", Value: cfg.GrpcCfg.OrderServiceGrpc},
		)
		return
	}
	defer conn.Close()
	o := pb.NewOrdersServiceClient(conn)
	orderGateway := order.NewGateway(o)

	orderUseCase := changedUc.NewOrderUseCase(orderGateway, deliveryUseCase, courierUseCase)

	handler := changedHandler.New(orderUseCase)

	saramaCfg := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(cfg.KafkaCfg.Version)
	if err != nil {
		lg.Fatal("invalid kafka version",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "version", Value: cfg.KafkaCfg.Version},
		)
	}
	saramaCfg.Version = version

	switch cfg.KafkaCfg.InitialOffset {
	case "oldest":
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		lg.Fatal("invalid initial offset",
			logger.Field{Key: "offset", Value: cfg.KafkaCfg.InitialOffset},
		)
	}

	saramaCfg.Consumer.Offsets.AutoCommit.Enable = cfg.KafkaCfg.AutoCommitEnable
	saramaCfg.Consumer.Offsets.AutoCommit.Interval = cfg.KafkaCfg.AutoCommitInterval

	kafkaClient, err := sarama.NewConsumerGroup(cfg.KafkaCfg.Brokers, cfg.KafkaCfg.GroupId, saramaCfg)
	if err != nil {
		lg.Fatal("unable to create kafka consumer group",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "brokers", Value: cfg.KafkaCfg.Brokers},
			logger.Field{Key: "groupId", Value: cfg.KafkaCfg.GroupId},
		)
		return
	}
	defer kafkaClient.Close()

	go func() {
		for {
			if err := kafkaClient.Consume(appCtx, []string{cfg.KafkaCfg.OrderTopic}, handler); err != nil {
				lg.Error("kafka consume error",
					logger.Field{Key: "error", Value: err},
					logger.Field{Key: "topic", Value: cfg.KafkaCfg.OrderTopic},
				)
			}

			if appCtx.Err() != nil {
				return
			}
		}
	}()

	<-shutdownCtx.Done()
	lg.Info("Shutting down kafka worker")

	_, gsCancel := context.WithTimeout(context.Background(), cfg.TimeCfg.ShutdownTimeout)
	defer gsCancel()

	lg.Info("Graceful shutdown kafka complete")
}
