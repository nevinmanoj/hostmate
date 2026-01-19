package payment

import (
	"time"
)

type PaymentType string

const (
	PaymentUPI          PaymentType = "UPI"
	PaymentCash         PaymentType = "Cash"
	PaymentBankTransfer PaymentType = "bank-transfer"
	PaymentOther        PaymentType = "other"
)

type Payment struct {
	ID          int64       `db:"id"`
	Amount      float64     `db:"amount"`
	Date        time.Time   `db:"date"`
	PaymentType PaymentType `db:"payment_type"`
	BookingID   int64       `db:"booking_id"`
	Remarks     string      `db:"remarks"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
	CreatedBy   int64       `db:"created_by"`
	UpdatedBy   int64       `db:"updated_by"`
}
