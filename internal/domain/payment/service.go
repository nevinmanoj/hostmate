package payment

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/nevinmanoj/hostmate/internal/domain/access"
	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
	property "github.com/nevinmanoj/hostmate/internal/domain/property"
	user "github.com/nevinmanoj/hostmate/internal/domain/user"
	middleware "github.com/nevinmanoj/hostmate/internal/middleware"
)

type PaymentService interface {
	GetAll(ctx context.Context, filter PaymentFilter) ([]Payment, int, error)
	GetWithBookingId(ctx context.Context, bookingID int64, page, pageSize int) ([]Payment, int, error)
	GetById(ctx context.Context, id int64) (*Payment, error)
	Create(ctx context.Context, property *Payment) error
	Update(ctx context.Context, property *Payment) error
}

type paymentService struct {
	repo          PaymentWriteRepository
	bookingRepo   booking.BookingReadRepository
	propertyRepo  property.PropertyReadRepository
	userRepo      user.UserReadRepository
	accessService access.AccessService
}

func NewPaymentService(
	repo PaymentWriteRepository,
	accessService access.AccessService,
	userRepo user.UserReadRepository,
	bookingRepo booking.BookingReadRepository,
	propertyRepo property.PropertyReadRepository) PaymentService {
	return &paymentService{
		repo:          repo,
		userRepo:      userRepo,
		bookingRepo:   bookingRepo,
		accessService: accessService,
		propertyRepo:  propertyRepo,
	}
}

func (s *paymentService) GetAll(ctx context.Context, filter PaymentFilter) ([]Payment, int, error) {

	userID := ctx.Value(middleware.ContextUserKey).(int64)
	// TODO setup admin bypass
	filter.UserID = &userID
	data, total, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		log.Println("Error fetching payments:", err)
		return nil, 0, ErrInternal
	}
	return data, total, nil
}
func (s *paymentService) GetWithBookingId(ctx context.Context, bookingID int64, limit, offset int) ([]Payment, int, error) {
	user := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessBooking(ctx, bookingID, user)
	if err != nil {
		return nil, 0, ErrInternal
	}
	if !hasAccess {
		return nil, 0, ErrUnauthorized
	}
	data, total, err := s.repo.GetByBookingId(ctx, bookingID, limit, offset)
	if err != nil {
		log.Printf("Error fetching payments for booking id %d:%s", bookingID, err.Error())
		return nil, 0, ErrInternal
	}
	return data, total, nil
}
func (s *paymentService) GetWithPropertyId(ctx context.Context, propertyID int64, limit, offset int) ([]Payment, int, error) {

	user := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessProperty(ctx, propertyID, user)
	if err != nil {
		return nil, 0, ErrInternal
	}
	if !hasAccess {
		return nil, 0, ErrUnauthorized
	}

	data, total, err := s.repo.GetByPropertyId(ctx, propertyID, limit, offset)
	if err != nil {
		log.Printf("Error fetching payments for property id %d:%s", propertyID, err.Error())
		return nil, 0, ErrInternal
	}
	return data, total, nil
}

func (s *paymentService) GetById(ctx context.Context, id int64) (*Payment, error) {
	user := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessPayment(ctx, id, user)
	if err != nil {
		return nil, ErrInternal
	}
	if !hasAccess {
		return nil, ErrUnauthorized
	}
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching payment with id %d: %s", id, err.Error())
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	return payment, nil
}

func (s *paymentService) Create(ctx context.Context, paymentToCreate *Payment) error {

	user := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessBooking(ctx, paymentToCreate.BookingID, user)
	if err != nil {
		return ErrInternal
	}
	if !hasAccess {
		return ErrUnauthorized
	}
	// Validate payment fields as needed, bookingID,images should exist
	createdBy, ok := ctx.Value(middleware.ContextUserKey).(int64)
	if !ok {
		log.Println("User not found in context")
		return ErrInternal
	}
	fmt.Printf("created by user:%d", createdBy)
	paymentToCreate.CreatedBy = createdBy
	paymentToCreate.UpdatedBy = createdBy
	err = s.repo.Create(ctx, paymentToCreate)
	if err != nil {
		log.Println("Error creating payment:", err)
		return ErrInternal
	}
	return nil
}

func (s *paymentService) Update(ctx context.Context, paymentToUpdate *Payment) error {
	user := ctx.Value(middleware.ContextUserKey).(int64)
	hasAccess, err := s.accessService.CanAccessPayment(ctx, paymentToUpdate.ID, user)
	if err != nil {
		return ErrInternal
	}
	if !hasAccess {
		return ErrUnauthorized
	}
	paymentFromDb, err := s.repo.GetByID(ctx, paymentToUpdate.ID)
	if err != nil {
		//no such property or unauthorized
		log.Printf("Error fetching payment with id %d: %s", paymentToUpdate.ID, err.Error())
		return err
	}

	paymentToUpdate.CreatedBy = paymentFromDb.CreatedBy
	paymentToUpdate.CreatedAt = paymentFromDb.CreatedAt
	paymentToUpdate.UpdatedBy = user
	err = s.repo.Update(ctx, paymentToUpdate)
	if err != nil {
		log.Println("Error updating property:", err)
		return ErrInternal
	}
	return nil
}
