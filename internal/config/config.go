package config

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

type Config struct {
	ServerCfg   *ServerConfig
	DatabaseCfg *DatabaseConfig
}

type ServerConfig struct {
	Port int `env:"PORT" env-default:"8080" validate:"required,gt=0,lt=65536"`
}

type TimeoutConfig struct {
	Shutdown time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
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

func MustLoad() *Config {
	var port int
	pflag.IntVar(&port, "port", 0, "HTTP server port")
	pflag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}
	cfg := &Config{
		ServerCfg:   &ServerConfig{},
		DatabaseCfg: &DatabaseConfig{},
	}

	if err := cleanenv.ReadEnv(cfg.ServerCfg); err != nil {
		log.Fatalf("Failed to load server config: %v", err)
	}

	if err := cleanenv.ReadEnv(cfg.DatabaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
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
