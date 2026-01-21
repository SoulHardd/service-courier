package config

import (
	"log"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

type Config struct {
	ServerCfg      *ServerConfig
	DatabaseCfg    *DatabaseConfig
	TimeCfg        *TimeConfig
	GrpcCfg        *GrpcConfig
	KafkaCfg       *KafkaConfig
	RateLimiterCfg *RateLimiterConfig
}

type ServerConfig struct {
	Port int `env:"PORT" env-default:"8080" validate:"required,gt=0,lt=65536"`
}

type TimeConfig struct {
	ShutdownTimeout         time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
	DeliveryMonitorInterval time.Duration `env:"DELIVERY_MONITOR_INTERVAL" env-default:"10s"`
	OrderWorkerInterval     time.Duration `env:"ORDER_WORKER_INTERVAL" env-default:"5s"`
}

type DatabaseConfig struct {
	Host            string        `env:"POSTGRES_HOST" validate:"required"`
	Port            string        `env:"POSTGRES_PORT" env-default:"5432" validate:"required"`
	User            string        `env:"POSTGRES_USER" validate:"required"`
	Password        string        `env:"POSTGRES_PASSWORD" validate:"required"`
	Dbname          string        `env:"POSTGRES_DB" validate:"required"`
	MaxConns        int32         `env:"POSTGRES_MAX_CONNS" env-default:"10" validate:"gt=0"`
	MinConns        int32         `env:"POSTGRES_MIN_CONNS" env-default:"5" validate:"gte=0"`
	MaxConnLifetime time.Duration `env:"POSTGRES_MAX_CONN_LIFETIME" env-default:"1h"`
	RetryMaxTime    time.Duration `env:"POSTGRES_RETRY_MAX_TIME" env-default:"30s"`
	RetryAttempts   int           `env:"POSTGRES_RETRY_ATTEMPTS" env-default:"5"`
}

type GrpcConfig struct {
	OrderServiceGrpc string        `env:"ORDER_SERVICE_GRPC"`
	MaxRetries       int           `env:"GRPC_MAX_RETRIES" env-default:"5"`
	BaseRetryDelay   time.Duration `env:"GRPC_BASE_RETRY_DELAY" env-default:"100ms"`
	DelayMultiplier  int           `env:"GRPC_DELAY_MULTIPLIER" env-default:"2"`
}

type KafkaConfig struct {
	Brokers            []string      `env:"KAFKA_BROKERS"`
	GroupId            string        `env:"KAFKA_GROUP_ID"`
	OrderTopic         string        `env:"KAFKA_ORDER_TOPIC"`
	Version            string        `env:"KAFKA_VERSION" env-default:"2.1.0"`
	InitialOffset      string        `env:"KAFKA_INITIAL_OFFSET" env-default:"oldest"`
	AutoCommitEnable   bool          `env:"KAFKA_AUTO_COMMIT_ENABLE" env-default:"true"`
	AutoCommitInterval time.Duration `env:"KAFKA_AUTO_COMMIT_INTERVAL" env-default:"1s"`
}

type RateLimiterConfig struct {
	MaxTokens  float64 `env:"RATE_LIMITER_MAX_TOKENS" env-default:"5"`
	RefillRate float64 `env:"RATE_LIMITER_REFILL_RATE" env-default:"5"`
}

func MustLoad() *Config {
	var port int
	pflag.IntVar(&port, "port", 0, "HTTP server port")
	pflag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}
	cfg := &Config{
		ServerCfg:      &ServerConfig{},
		DatabaseCfg:    &DatabaseConfig{},
		TimeCfg:        &TimeConfig{},
		GrpcCfg:        &GrpcConfig{},
		KafkaCfg:       &KafkaConfig{},
		RateLimiterCfg: &RateLimiterConfig{},
	}

	if err := cleanenv.ReadEnv(cfg.ServerCfg); err != nil {
		log.Fatalf("Failed to load server config: %v", err)
	}
	if err := cleanenv.ReadEnv(cfg.DatabaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}
	if err := cleanenv.ReadEnv(cfg.TimeCfg); err != nil {
		log.Fatalf("Failed to load time config: %v", err)
	}
	if err := cleanenv.ReadEnv(cfg.GrpcCfg); err != nil {
		log.Fatalf("Failed to load grpc config: %v", err)
	}
	if err := cleanenv.ReadEnv(cfg.KafkaCfg); err != nil {
		log.Fatalf("Failed to load kafka config: %v", err)
	}
	if err := cleanenv.ReadEnv(cfg.RateLimiterCfg); err != nil {
		log.Fatalf("Failed to load rate limiter config: %v", err)
	}

	if port != 0 {
		cfg.ServerCfg.Port = port
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	return cfg
}
