package payment

import (
	"time"

	payment "github.com/nevinmanoj/hostmate/internal/domain/payment"
)

type CreatePaymentRequest struct {
	Amount      float64             `json:"amount"`
	Date        time.Time           `json:"date"`
	PaymentType payment.PaymentType `json:"payment_type"`
	BookingID   int64               `json:"booking_id"`
	Remarks     string              `json:"remarks"`
}

type UpdatePaymentRequest struct {
	ID int64 `json:"id"`
	CreatePaymentRequest
}

type PaymentResponse struct {
	ID          int64               `json:"id"`
	Amount      float64             `json:"amount"`
	Date        time.Time           `json:"date"`
	PaymentType payment.PaymentType `json:"payment_type"`
	BookingID   int64               `json:"booking_id"`
	Remarks     string              `json:"remarks"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	CreatedBy   int64               `json:"created_by"`
	UpdatedBy   int64               `json:"updated_by"`
}

func ToPaymentResponse(p *payment.Payment) PaymentResponse {
	return PaymentResponse{
		ID:          p.ID,
		Amount:      p.Amount,
		Date:        p.Date,
		PaymentType: p.PaymentType,
		BookingID:   p.BookingID,
		Remarks:     p.Remarks,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		CreatedBy:   p.CreatedBy,
		UpdatedBy:   p.UpdatedBy,
	}
}
