package domain

import "errors"

var (
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrPhoneExists           = errors.New("courier with this phone already exists")
	ErrCourierNotFound       = errors.New("courier not found")
	ErrInvalidPhone          = errors.New("invalid phone number")
)
