package httperror

import (
	"avito/internal/domain"

	"github.com/pkg/errors"
)

type HTTPError struct {
	Code    int
	Message string
}

var (
	ErrInvalidJSON      = HTTPError{Code: 400, Message: `{"error": "Invalid JSON"}`}
	ErrInvalidCourierId = HTTPError{Code: 400, Message: `{"error": "Invalid courier id"}`}
	ErrDB               = HTTPError{Code: 500, Message: `{"error": "Database error"}`}
)

func MapDomainError(err error) HTTPError {
	switch {
	case errors.Is(err, domain.ErrNoAvailableCouriers):
		return HTTPError{Code: 409, Message: `{"error": "No available couriers"}`}
	case errors.Is(err, domain.ErrDeliveryAlreadyExists):
		return HTTPError{Code: 409, Message: `{"error": "Delivery already exists"}`}
	case errors.Is(err, domain.ErrDeliveryNotFound):
		return HTTPError{Code: 404, Message: `{"error": "Delivery not found"}`}

	case errors.Is(err, domain.ErrCourierNotFound):
		return HTTPError{Code: 404, Message: `{"error": "Courier not found"}`}
	case errors.Is(err, domain.ErrPhoneExists):
		return HTTPError{Code: 409, Message: `{"error": "Courier with this phone already exists"}`}
	case errors.Is(err, domain.ErrInvalidPhone):
		return HTTPError{Code: 400, Message: `{"error": "Invalid phone"}`}

	case errors.Is(err, domain.ErrMissingRequiredFields):
		return HTTPError{Code: 400, Message: `{"error": "Missing required fields"}`}
	case errors.Is(err, domain.ErrInvalidId):
		return HTTPError{Code: 400, Message: `{"error": "Invalid ID"}`}

	default:
		return ErrDB
	}
}
