package courier

import (
	"avito/internal/domain"
	"avito/internal/model"
	"context"
	"errors"
	"fmt"

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

func (r *CourierRepository) GetOneById(ctx context.Context, id int64) (*model.Courier, error) {
	var courier model.CourierDB

	err := r.pool.QueryRow(ctx,
		"SELECT id, name, phone, status FROM couriers WHERE id=$1",
		id).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Phone,
		&courier.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrCourierNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.Courier{
		ID:     courier.ID,
		Name:   courier.Name,
		Phone:  courier.Phone,
		Status: courier.Status,
	}, nil
}

func (r *CourierRepository) GetAll(ctx context.Context) ([]model.Courier, error) {

	rows, err := r.pool.Query(ctx, "SELECT id, name, phone, status FROM couriers ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var couriers []model.Courier

	for rows.Next() {
		var courier model.CourierDB

		err := rows.Scan(
			&courier.ID,
			&courier.Name,
			&courier.Phone,
			&courier.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
		couriers = append(couriers, model.Courier{
			ID:     courier.ID,
			Name:   courier.Name,
			Phone:  courier.Phone,
			Status: courier.Status,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return couriers, nil
}

func (r *CourierRepository) Create(ctx context.Context, courier *model.Courier) (int64, error) {
	var id int64

	courierDB := model.CourierDB{
		Name:   courier.Name,
		Phone:  courier.Phone,
		Status: courier.Status,
	}

	err := r.pool.QueryRow(ctx,
		"INSERT INTO couriers (name, phone, status) VALUES ($1, $2, $3) RETURNING id",
		courierDB.Name, courierDB.Phone, courierDB.Status).Scan(&id)

	const duplicateKeyCode = "23505"
	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == duplicateKeyCode {
			return 0, domain.ErrPhoneExists
		}
		return 0, fmt.Errorf("database error: %w", err)
	}
	return id, nil
}

func (r *CourierRepository) Update(ctx context.Context, courier *model.Courier) error {
	courierDB := model.CourierDB{
		ID:     courier.ID,
		Name:   courier.Name,
		Phone:  courier.Phone,
		Status: courier.Status,
	}

	result, err := r.pool.Exec(ctx, "UPDATE couriers SET name=$1, phone=$2, status=$3, updated_at=NOW() WHERE id=$4",
		courierDB.Name, courierDB.Phone, courierDB.Status, courierDB.ID)

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
