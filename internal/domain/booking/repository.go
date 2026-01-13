package booking

import (
	"context"
	"time"
)

type BookingReadRepository interface {
	GetAll(ctx context.Context, filter BookingFilter) ([]Booking, int, error)
	GetByID(ctx context.Context, id int64) (*Booking, error)
	CheckAvailability(ctx context.Context, propertyID int64, checkInDate, checkOutDate time.Time) (bool, error)
}
type BookingWriteRepository interface {
	BookingReadRepository
	Create(ctx context.Context, property *Booking) error
	Update(ctx context.Context, property *Booking) error
}
