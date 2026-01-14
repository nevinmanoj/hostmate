package payment

import (
	"time"
)

type PaymentFilter struct {
	UserID      *int64
	FromDate    *time.Time
	ToDate      *time.Time
	PaymentType []PaymentType
	Limit       int
	Offset      int
}
