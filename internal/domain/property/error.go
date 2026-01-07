package property

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("Not found")
	ErrNotValidManagers = errors.New("managers are not valid")
	ErrInternal         = errors.New("internal error")
	ErrUnauthorized     = errors.New("unauthorized")
)
