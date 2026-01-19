package access

import (
	"context"

	domaincore "github.com/nevinmanoj/hostmate/internal/domain/core"
)

type AccessService interface {
	CanAccessPayment(ctx context.Context, paymentID, userID int64) (bool, error)
	CanAccessBooking(ctx context.Context, bookingID, userID int64) (bool, error)
	CanAccessProperty(ctx context.Context, propertyID, userID int64) (bool, error)
	CanAccessAttachment(ctx context.Context, attachmentParentType domaincore.AttachmentParentType, attachemntID, userID int64) (bool, error)
}

type accessService struct {
	repo AccessRepository
}

func NewAccessService(repo AccessRepository) AccessService {
	return &accessService{repo: repo}
}

func (s *accessService) CanAccessPayment(ctx context.Context, paymentID, userID int64) (bool, error) {
	canAccess, err := s.repo.HasManagerByPaymentID(ctx, paymentID, userID)
	if err != nil {
		return false, err
	}
	return canAccess, nil
}
func (s *accessService) CanAccessBooking(ctx context.Context, bookingID, userID int64) (bool, error) {
	canAccess, err := s.repo.HasManagerByBookingID(ctx, bookingID, userID)
	if err != nil {
		return false, err
	}
	return canAccess, nil
}
func (s *accessService) CanAccessProperty(ctx context.Context, propertyID, userID int64) (bool, error) {
	canAccess, err := s.repo.HasManagerByPropertyID(ctx, propertyID, userID)
	if err != nil {
		return false, err
	}
	return canAccess, nil
}
func (s *accessService) CanAccessAttachment(ctx context.Context, attachmentParentType domaincore.AttachmentParentType, attachemntID, userID int64) (bool, error) {
	canAccess, err := s.repo.HasManagerByAttachmentID(ctx, attachmentParentType, attachemntID, userID)
	if err != nil {
		return false, err
	}
	return canAccess, nil
}
