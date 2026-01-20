package courier

import (
	"avito/internal/domain"
	"context"
	"errors"
	"regexp"
)

type CourierUseCase struct {
	repository CourierRepository
}

func New(repository CourierRepository) *CourierUseCase {
	return &CourierUseCase{repository: repository}
}

func (u *CourierUseCase) GetCourier(ctx context.Context, id int64) (*domain.Courier, error) {
	courier, err := u.repository.GetOneById(ctx, id)

	if err != nil {
		if errors.Is(err, domain.ErrCourierNotFound) {
			return nil, domain.ErrCourierNotFound
		}
		return nil, err
	}

	return courier, nil
}

func (u *CourierUseCase) GetCouriers(ctx context.Context) ([]domain.Courier, error) {
	return u.repository.GetAll(ctx)
}

func (u *CourierUseCase) CreateCourier(ctx context.Context, req *domain.Courier) (int64, error) {
	if req.Name == "" || req.Phone == "" || req.Status == "" {
		return 0, domain.ErrMissingRequiredFields
	}

	if isValid := validatePhoneNumber(req.Phone); !isValid {
		return 0, domain.ErrInvalidPhone
	}

	id, err := u.repository.Create(ctx, req)

	if err != nil {
		if errors.Is(err, domain.ErrPhoneExists) {
			return 0, domain.ErrPhoneExists
		}
		return 0, err
	}

	return id, nil
}

func (u *CourierUseCase) UpdateCourier(ctx context.Context, req *domain.Courier) error {
	if req.Name == "" || req.Phone == "" || req.Status == "" {
		return domain.ErrMissingRequiredFields
	}

	if isValid := validatePhoneNumber(req.Phone); !isValid {
		return domain.ErrInvalidPhone
	}

	err := u.repository.Update(ctx, req)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCourierNotFound):
			return domain.ErrCourierNotFound
		case errors.Is(err, domain.ErrPhoneExists):
			return domain.ErrPhoneExists
		}
		return err
	}

	return nil
}

func (u *CourierUseCase) ReleaseCourierByOrderId(ctx context.Context, orderId string) error {
	err := u.repository.ReleaseOneByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, domain.ErrDeliveryNotFound) {
			return domain.ErrDeliveryNotFound
		}
		return err
	}
	return nil
}

var phoneRegex = regexp.MustCompile(`^\+\d{11}$`)

func validatePhoneNumber(phone string) bool {
	return phoneRegex.MatchString(phone)
}
