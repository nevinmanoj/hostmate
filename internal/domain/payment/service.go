package payment

import (
	"context"
	"errors"
	"fmt"
	"log"

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
	repo         PaymentWriteRepository
	bookingRepo  booking.BookingReadRepository
	propertyRepo property.PropertyReadRepository
	userRepo     user.UserReadRepository
}

func NewPaymentService(
	repo PaymentWriteRepository,
	userRepo user.UserReadRepository,
	bookingRepo booking.BookingReadRepository,
	propertyRepo property.PropertyReadRepository) PaymentService {
	return &paymentService{repo: repo, userRepo: userRepo, bookingRepo: bookingRepo}
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

	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		log.Printf("Invalid booking id %d: %s", bookingID, err.Error())
		return nil, 0, ErrNotValidBookingId
	}
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	ok, err := s.propertyRepo.HasManager(ctx, booking.PropertyID, userID)
	if err != nil {
		log.Println("error checking managers", err.Error())
		return nil, 0, ErrInternal
	}
	if !ok {
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

	userID := ctx.Value(middleware.ContextUserKey).(int64)
	ok, err := s.propertyRepo.HasManager(ctx, propertyID, userID)
	if err != nil {
		log.Println("error checking managers", err.Error())
		return nil, 0, ErrInternal
	}
	if !ok {
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
	payment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching payment with id %d: %s", id, err.Error())
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}
	bookingID := payment.BookingID
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		log.Printf("error fetching booking with for payment id %d: %s", bookingID, err.Error())
		return nil, ErrInternal
	}
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	ok, err := s.propertyRepo.HasManager(ctx, booking.PropertyID, userID)
	if err != nil {
		log.Println("error checking managers", err.Error())
		return nil, ErrInternal
	}
	if !ok {
		return nil, ErrUnauthorized
	}
	return payment, nil
}

func (s *paymentService) Create(ctx context.Context, paymentToCreate *Payment) error {

	bookingID := paymentToCreate.BookingID
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		log.Printf("error fetching booking with for payment id %d: %s", bookingID, err.Error())
		return ErrInternal
	}
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	ok, err := s.propertyRepo.HasManager(ctx, booking.PropertyID, userID)
	if err != nil {
		log.Println("error checking managers", err.Error())
		return ErrInternal
	}
	if !ok {
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
	// Validate property fields as needed, bookingID,images should exist
	paymentFromDb, err := s.repo.GetByID(ctx, paymentToUpdate.ID)
	if err != nil {
		//no such property or unauthorized
		log.Printf("Error fetching payment with id %d: %s", paymentToUpdate.ID, err.Error())
		return err
	}
	user, ok := ctx.Value(middleware.ContextUserKey).(int64)
	if !ok {
		log.Println("User not found in context")
		return ErrInternal
	}
	bookingID := paymentToUpdate.BookingID
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		log.Printf("error fetching booking with for payment id %d: %s", bookingID, err.Error())
		return ErrInternal
	}
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	ok, err = s.propertyRepo.HasManager(ctx, booking.PropertyID, userID)
	if err != nil {
		log.Println("error checking managers", err.Error())
		return ErrInternal
	}
	if !ok {
		return ErrUnauthorized
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
