package access

import (
	"context"

	domaincore "github.com/nevinmanoj/hostmate/internal/domain/core"
)

type AccessRepository interface {
	HasManagerByPropertyID(ctx context.Context, propertyID, userID int64) (bool, error)
	HasManagerByBookingID(ctx context.Context, bookingID, userID int64) (bool, error)
	HasManagerByPaymentID(ctx context.Context, paymentID, userID int64) (bool, error)
	HasManagerByAttachmentID(ctx context.Context, parentType domaincore.AttachmentParentType, AttachmentID, userID int64) (bool, error)
}
