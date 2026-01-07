package booking

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("Not found")
	ErrInternal         = errors.New("internal error")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInvalidDateRange = errors.New("invalid date range")
	ErrBookingConflict  = errors.New("booking conflict")
)
