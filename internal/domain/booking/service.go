package booking

import (
	"context"
	"time"

	"github.com/nevinmanoj/hostmate/internal/domain/access"
	"github.com/nevinmanoj/hostmate/internal/domain/property"
	"github.com/nevinmanoj/hostmate/internal/middleware"
)

type BookingService interface {
	GetAll(ctx context.Context, filter BookingFilter) ([]Booking, int, error)
	GetById(ctx context.Context, id int64) (*Booking, error)
	Create(ctx context.Context, booking *Booking) error
	Update(ctx context.Context, booking *Booking) error
	CheckAvailability(ctx context.Context, propertyID int64, startDate, endDate time.Time) (bool, error)
}

type bookingService struct {
	repo          BookingWriteRepository
	propertyRepo  property.PropertyReadRepository
	accessService access.AccessService
}

func NewBookingService(repo BookingWriteRepository, propertyRepo property.PropertyReadRepository, accessService access.AccessService) BookingService {
	return &bookingService{repo: repo, propertyRepo: propertyRepo, accessService: accessService}
}
func (s *bookingService) GetAll(ctx context.Context, filter BookingFilter) ([]Booking, int, error) {
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	// TODO setup admin bypass
	filter.UserID = &userID
	data, total, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return data, total, nil
}

func (s *bookingService) GetById(ctx context.Context, id int64) (*Booking, error) {
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessBooking(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, ErrUnauthorized
	}
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) Create(ctx context.Context, booking *Booking) error {
	// Validate property fields as needed, managers
	createdBy, ok := ctx.Value(middleware.ContextUserKey).(int64)
	if !ok {
		return ErrInternal
	}
	hasAccess, err := s.accessService.CanAccessProperty(ctx, booking.PropertyID, createdBy)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrUnauthorized
	}

	//check if booking dates are valid
	if !booking.CheckInDate.Before(booking.CheckOutDate) {
		return ErrInvalidDateRange
	}
	//check property availability for the booking dates can be added here
	booking.CreatedBy = createdBy
	booking.UpdatedBy = createdBy
	booking.ManagerID = createdBy
	err = s.repo.Create(ctx, booking)
	if err != nil {
		return err
	}
	return nil
}

func (s *bookingService) Update(ctx context.Context, booking *Booking) error {
	// Validate booking fields as needed, managers,images should exist
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessBooking(ctx, booking.ID, userID)
	if err != nil {
		return err
	}
	if !hasAccess {
		return ErrUnauthorized
	}

	bookingFromDb, err := s.repo.GetByID(ctx, booking.ID)
	if err != nil {
		//no such booking
		return err
	}

	//check if booking dates are valid
	if !booking.CheckInDate.Before(booking.CheckOutDate) {
		return ErrInvalidDateRange
	}
	booking.CreatedBy = bookingFromDb.CreatedBy
	booking.CreatedAt = bookingFromDb.CreatedAt
	booking.UpdatedBy = userID
	err = s.repo.Update(ctx, booking)
	if err != nil {
		return err
	}
	return nil
}

func (s *bookingService) CheckAvailability(ctx context.Context, propertyID int64, startDate, endDate time.Time) (bool, error) {
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessProperty(ctx, propertyID, userID)
	if err != nil {
		return false, err
	}
	if !hasAccess {
		return false, ErrUnauthorized
	}
	available, err := s.repo.CheckAvailability(ctx, propertyID, startDate, endDate)

	return available, nil
}
