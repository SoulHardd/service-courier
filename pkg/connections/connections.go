package connections

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/SoulHardd/service-courier/internal/config"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitPool(ctx context.Context, dbCfg *config.DatabaseConfig) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Dbname))

	if err != nil {
		log.Fatal(err)
	}

	cfg.MaxConns = dbCfg.MaxConns
	cfg.MaxConnLifetime = dbCfg.MaxConnLifetime
	cfg.MinConns = dbCfg.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = PingWithRetry(ctx, pool, dbCfg.RetryAttempts, dbCfg.RetryMaxTime)
	if err != nil {
		log.Fatal(err)
	}

	return pool
}

func PingWithRetry(ctx context.Context, pool *pgxpool.Pool, retries int, maxTime time.Duration) error {
	backoffCfg := backoff.NewExponentialBackOff()
	backoffCfg.MaxElapsedTime = maxTime

	operation := func() error {
		return pool.Ping(ctx)
	}

	return backoff.Retry(operation, backoff.WithMaxRetries(backoffCfg, uint64(retries)))
}
