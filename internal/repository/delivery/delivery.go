package delivery

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SoulHardd/service-courier/internal/domain"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *DeliveryRepository {
	return &DeliveryRepository{
		pool: pool,
	}
}

func (r *DeliveryRepository) Assign(ctx context.Context, orderId string, c *domain.Courier, deadline time.Time) (*domain.Delivery, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		err = errors.Join(tx.Rollback(ctx))
	}()

	if err := validateUuid(orderId); err != nil {
		if errors.Is(err, domain.ErrInvalidId) {
			return nil, domain.ErrInvalidId
		}
		return nil, err
	}

	courierQuery, courierArgs, _ := squirrel.
		Select("status").
		From("couriers").
		Where(squirrel.Eq{"id": c.ID}).
		Suffix("FOR UPDATE").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var courierStatus string

	if err := tx.QueryRow(ctx, courierQuery, courierArgs...).Scan(&courierStatus); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCourierNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	if domain.CourierStatus(courierStatus) == domain.CourierStatusBusy {
		return nil, domain.ErrCourierIsBusy
	}

	updCourierQuery, updCourierArgs, _ := squirrel.
		Update("couriers").
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("status", string(domain.CourierStatusBusy)).
		Where(squirrel.Eq{"id": c.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := tx.Exec(ctx, updCourierQuery, updCourierArgs...); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	delivery := DeliveryDB{
		CourierId: c.ID,
		OrderId:   orderId,
		Deadline:  deadline,
	}

	checkDeliveryQuery, checkDeliveryArgs, _ := squirrel.
		Select("COUNT(*)").
		From("delivery").
		Where(squirrel.Eq{"order_id": orderId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var existingDeliveries int
	if err := tx.QueryRow(ctx, checkDeliveryQuery, checkDeliveryArgs...).Scan(&existingDeliveries); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if existingDeliveries > 0 {
		return nil, domain.ErrDeliveryAlreadyExists
	}

	deliveryQuery, deliveryArgs, _ := squirrel.
		Insert("delivery").
		Columns("courier_id", "order_id", "deadline").
		Values(delivery.CourierId, delivery.OrderId, delivery.Deadline).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err := tx.QueryRow(ctx, deliveryQuery, deliveryArgs...).Scan(&delivery.ID); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}
	return &domain.Delivery{
		CourierId:     delivery.CourierId,
		OrderId:       delivery.OrderId,
		TransportType: c.TransportType,
		Deadline:      delivery.Deadline,
	}, nil

}

func (r *DeliveryRepository) Unassign(ctx context.Context, orderId string) (*domain.Delivery, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		err = errors.Join(tx.Rollback(ctx))
	}()

	if err := validateUuid(orderId); err != nil {
		if errors.Is(err, domain.ErrInvalidId) {
			return nil, domain.ErrInvalidId
		}
		return nil, err
	}

	selectQuery, selectArgs, _ := squirrel.
		Select("courier_id", "order_id").
		From("delivery").
		Where(squirrel.Eq{"order_id": orderId}).
		Suffix("FOR UPDATE").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var delivery DeliveryDB

	if err := tx.QueryRow(ctx, selectQuery, selectArgs...).Scan(&delivery.CourierId, &delivery.OrderId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDeliveryNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	updCourierQuery, updCourierArgs, _ := squirrel.
		Update("couriers").
		Set("updated_at", squirrel.Expr("NOW()")).
		Set("status", string(domain.CourierStatusAvailable)).
		Where(squirrel.Eq{"id": delivery.CourierId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := tx.Exec(ctx, updCourierQuery, updCourierArgs...); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	deleteQuery, deleteArgs, _ := squirrel.
		Delete("delivery").
		Where(squirrel.Eq{"order_id": orderId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := tx.Exec(ctx, deleteQuery, deleteArgs...); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &domain.Delivery{
		OrderId:   delivery.OrderId,
		CourierId: delivery.CourierId,
	}, nil
}

func (r *DeliveryRepository) ReleaseExpireDeliveries(ctx context.Context) error {

	query, args, _ := squirrel.
		Update("couriers").
		Set("status", string(domain.CourierStatusAvailable)).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where("NOT EXISTS (SELECT 1 FROM delivery WHERE courier_id = couriers.id AND deadline >= NOW())").
		Where(squirrel.Eq{"status": string(domain.CourierStatusBusy)}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if _, err := r.pool.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("database error: %w", err)
	}

	return nil
}

func validateUuid(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return domain.ErrInvalidId
	}
	return nil
}
