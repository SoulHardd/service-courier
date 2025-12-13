//go:build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"avito/internal/config"
	courierRepo "avito/internal/repository/courier"
	deliveryRepo "avito/internal/repository/delivery"
	courierUC "avito/internal/useCase/courier"
	deliveryUC "avito/internal/useCase/delivery"
	"avito/pkg/connections"

	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDependencies struct {
	Pool            *pgxpool.Pool
	CourierRepo     *courierRepo.CourierRepository
	DeliveryRepo    *deliveryRepo.DeliveryRepository
	CourierUseCase  *courierUC.CourierUseCase
	DeliveryUseCase *deliveryUC.DeliveryUseCase
	Container       testcontainers.Container
}

func SetupTestDB(ctx context.Context, t *testing.T) (testcontainers.Container, *pgxpool.Pool, error) {
	t.Helper()

	container, err := postgres.RunContainer(ctx,
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),

		testcontainers.WithImage("postgres:15"),

		postgres.WithInitScripts(),

		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get container connection string: %w", err)
	}

	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolCfg.MaxConns = 5
	poolCfg.MinConns = 1
	poolCfg.MaxConnLifetime = time.Minute * 10

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	dbCfg := &config.DatabaseConfig{
		RetryAttempts: 3,
		RetryMaxTime:  5 * time.Second,
	}
	if err := connections.PingWithRetry(ctx, pool, dbCfg.RetryAttempts, dbCfg.RetryMaxTime); err != nil {
		pool.Close()
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := initTestSchema(ctx, pool); err != nil {
		pool.Close()
		_ = container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to init schema: %w", err)
	}

	return container, pool, nil
}

func initTestSchema(ctx context.Context, pool *pgxpool.Pool) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS couriers (
			 id             BIGSERIAL PRIMARY KEY,
			 name           TEXT NOT NULL,
			 phone          TEXT NOT NULL UNIQUE,
			 status         TEXT NOT NULL DEFAULT 'available',
			 transport_type TEXT NOT NULL DEFAULT 'on_foot',
			 created_at     TIMESTAMP DEFAULT now(),
			 updated_at     TIMESTAMP DEFAULT now()
		)`,
		`CREATE TABLE IF NOT EXISTS delivery (
			 id          BIGSERIAL PRIMARY KEY,
			 courier_id  BIGINT NOT NULL,
			 order_id    VARCHAR(255) NOT NULL,
			 assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
			 deadline    TIMESTAMP NOT NULL
		)`,
	}

	for _, q := range queries {
		if _, err := pool.Exec(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
func SetupTestDependencies(ctx context.Context, t *testing.T) (*TestDependencies, error) {
	t.Helper()

	container, pool, err := SetupTestDB(ctx, t)
	if err != nil {
		return nil, err
	}

	cRepo := courierRepo.New(pool)
	dRepo := deliveryRepo.New(pool)

	cUC := courierUC.New(cRepo)
	dUC := deliveryUC.New(dRepo, cRepo)

	deps := &TestDependencies{
		Pool:            pool,
		CourierRepo:     cRepo,
		DeliveryRepo:    dRepo,
		CourierUseCase:  cUC,
		DeliveryUseCase: dUC,
		Container:       container,
	}
	return deps, nil
}

func CleanTestData(ctx context.Context, pool *pgxpool.Pool) error {
	queries := []string{
		"TRUNCATE TABLE delivery RESTART IDENTITY CASCADE",
		"TRUNCATE TABLE couriers RESTART IDENTITY CASCADE",
	}

	for _, q := range queries {
		if _, err := pool.Exec(ctx, q); err != nil {
			return err
		}
	}
	return nil
}
