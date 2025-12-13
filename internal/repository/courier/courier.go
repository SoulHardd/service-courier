package courier

import (
	"avito/internal/domain"
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *CourierRepository {
	return &CourierRepository{
		pool: pool,
	}
}

func (r *CourierRepository) GetOneById(ctx context.Context, id int64) (*domain.Courier, error) {
	var courier CourierDB

	query, args, _ := squirrel.
		Select("id", "name", "phone", "status", "transport_type").
		From("couriers").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Phone,
		&courier.Status,
		&courier.TransportType,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCourierNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &domain.Courier{
		ID:            courier.ID,
		Name:          courier.Name,
		Phone:         courier.Phone,
		Status:        courier.Status,
		TransportType: courier.TransportType,
	}, nil
}

func (r *CourierRepository) GetAll(ctx context.Context) ([]domain.Courier, error) {

	query, args, _ := squirrel.
		Select("id", "name", "status", "transport_type").
		From("couriers").
		OrderBy("id ASC").
		ToSql()

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var couriers []domain.Courier

	for rows.Next() {
		var courier CourierDB

		err := rows.Scan(
			&courier.ID,
			&courier.Name,
			&courier.Status,
			&courier.TransportType,
		)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		couriers = append(couriers, domain.Courier{
			ID:            courier.ID,
			Name:          courier.Name,
			Status:        courier.Status,
			TransportType: courier.TransportType,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return couriers, nil
}

func (r *CourierRepository) Create(ctx context.Context, courier *domain.Courier) (int64, error) {
	var id int64

	courierDB := CourierDB{
		Name:          courier.Name,
		Phone:         courier.Phone,
		Status:        courier.Status,
		TransportType: courier.TransportType,
	}

	query, args, _ := squirrel.
		Insert("couriers").
		Columns("name", "phone", "status", "transport_type").
		Values(courierDB.Name, courierDB.Phone, courierDB.Status, courierDB.TransportType).
		Suffix("RETURNING id").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	err := r.pool.QueryRow(ctx, query, args...).Scan(&id)

	const duplicateKeyCode = "23505"
	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == duplicateKeyCode {
			return 0, domain.ErrPhoneExists
		}
		return 0, fmt.Errorf("database error: %w", err)
	}
	return id, nil
}

func (r *CourierRepository) Update(ctx context.Context, courier *domain.Courier) error {
	courierDB := CourierDB{
		ID:            courier.ID,
		Name:          courier.Name,
		Phone:         courier.Phone,
		Status:        courier.Status,
		TransportType: courier.TransportType,
	}

	query, args, _ := squirrel.
		Update("couriers").
		Set("name", courierDB.Name).
		Set("phone", courier.Phone).
		Set("status", courierDB.Status).
		Set("transport_type", courier.TransportType).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": courier.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	result, err := r.pool.Exec(ctx, query, args...)

	const duplicateKeyCode = "23505"
	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == duplicateKeyCode {
			return domain.ErrPhoneExists
		}
		return fmt.Errorf("database error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrCourierNotFound
	}
	return nil
}

func (r *CourierRepository) GetOneForDelivery(ctx context.Context) (*domain.Courier, error) {

	query, args, _ := squirrel.
		Select("c.id", "c.name", "c.phone", "c.status", "c.transport_type").
		From("couriers c").
		LeftJoin("delivery d ON c.id = d.courier_id AND d.deadline < NOW()").
		Where(squirrel.Eq{"c.status": string(domain.CourierStatusAvailable)}).
		GroupBy("c.id", "c.name", "c.phone", "c.status", "c.transport_type").
		OrderBy("COUNT(d.id) ASC", "c.id ASC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	var courier CourierDB
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Phone,
		&courier.Status,
		&courier.TransportType,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNoAvailableCouriers
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &domain.Courier{
		ID:            courier.ID,
		Name:          courier.Name,
		Phone:         courier.Phone,
		Status:        courier.Status,
		TransportType: courier.TransportType,
	}, nil
}
