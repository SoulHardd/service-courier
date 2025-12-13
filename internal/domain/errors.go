package domain

import "errors"

var (
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrPhoneExists           = errors.New("courier with this phone already exists")
	ErrCourierNotFound       = errors.New("courier not found")
	ErrInvalidPhone          = errors.New("invalid phone number")
	ErrNoAvailableCouriers   = errors.New("no available couriers")
	ErrUnsupportedTransport  = errors.New("unsupported transport type")
	ErrCourierIsBusy         = errors.New("courier is busy")
	ErrDeliveryNotFound      = errors.New("delivery not found")
	ErrInvalidId             = errors.New("invalid id")
	ErrDeliveryAlreadyExists = errors.New("delivery already exists")
)
