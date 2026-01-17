package access

import "context"

type AccessRepository interface {
	HasManagerByPropertyID(ctx context.Context, propertyID, userID int64) (bool, error)
	HasManagerByBookingID(ctx context.Context, bookingID, userID int64) (bool, error)
	HasManagerByPaymentID(ctx context.Context, paymentID, userID int64) (bool, error)
}
