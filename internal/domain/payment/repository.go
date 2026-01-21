package payment

import (
	"context"
)

type PaymentReadRepository interface {
	GetAll(ctx context.Context, filter PaymentFilter) ([]Payment, int, error)
	GetByBookingId(ctx context.Context, bookingID int64, limit, offset int) ([]Payment, int, error)
	GetByPropertyId(ctx context.Context, propertyID int64, limit, offset int) ([]Payment, int, error)
	GetByID(ctx context.Context, id int64) (*Payment, error)
	GetBlobs(ctx context.Context, paymentID int64) ([]string, error)
}
type PaymentWriteRepository interface {
	PaymentReadRepository
	Create(ctx context.Context, property *Payment) error
	Update(ctx context.Context, property *Payment) error
	AppendBlobs(ctx context.Context, paymentID int64, blobName string) error
}
