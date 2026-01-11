package payment

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("Payment not found")
	ErrNotValidBookingId = errors.New("Booking Id is not valid")
	ErrInternal          = errors.New("Internal error")
	ErrUnauthorized      = errors.New("unauthorized")
)
