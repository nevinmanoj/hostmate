package payment

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"

	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
	user "github.com/nevinmanoj/hostmate/internal/domain/user"
	middleware "github.com/nevinmanoj/hostmate/internal/middleware"
)

type PaymentService interface {
	GetAll(ctx context.Context, page, pageSize int) ([]Payment, int, error)
	GetWithBookingId(ctx context.Context, bookingID int64, page, pageSize int) ([]Payment, int, error)
	GetById(ctx context.Context, id int64) (*Payment, error)
	Create(ctx context.Context, property *Payment) error
	Update(ctx context.Context, property *Payment) error
}

type paymentService struct {
	repo        PaymentWriteRepository
	bookingRepo booking.BookingReadRepository
	userRepo    user.UserReadRepository
}

func NewPaymentService(repo PaymentWriteRepository, userRepo user.UserReadRepository, bookingRepo booking.BookingReadRepository) PaymentService {
	return &paymentService{repo: repo, userRepo: userRepo, bookingRepo: bookingRepo}
}

func (s *paymentService) GetAll(ctx context.Context, page, pageSize int) ([]Payment, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	data, total, err := s.repo.GetAll(ctx, pageSize, offset)
	if err != nil {
		log.Println("Error fetching payments:", err)
		return nil, 0, ErrInternal
	}
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	return data, totalPages, nil
}
func (s *paymentService) GetWithBookingId(ctx context.Context, bookingID int64, page, pageSize int) ([]Payment, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		log.Printf("Invalid booking id %d: %s", bookingID, err.Error())
		return nil, 0, ErrNotValidBookingId
	}

	data, total, err := s.repo.GetByBookingId(ctx, bookingID, pageSize, offset)
	if err != nil {
		log.Printf("Error fetching payments for booking id %d:%s", bookingID, err.Error())
		return nil, 0, ErrInternal
	}
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	return data, totalPages, nil
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
		log.Printf("Invalid/Unauthorized booking id %d: %s", bookingID, err.Error())
		return nil, ErrUnauthorized
	}

	return payment, nil
}

func (s *paymentService) Create(ctx context.Context, paymentToCreate *Payment) error {
	// Validate payment fields as needed, bookingID,images should exist
	bookingID := paymentToCreate.BookingID
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil || booking == nil {
		log.Printf("Invalid/Unauthorized booking id %d: %s", bookingID, err.Error())
		return ErrUnauthorized
	}
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
