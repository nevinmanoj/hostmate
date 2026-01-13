package payment

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	payment "github.com/nevinmanoj/hostmate/internal/domain/payment"
)

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentReadRepository(db *sqlx.DB) payment.PaymentReadRepository {
	return &paymentRepository{db: db}
}
func NewPaymentWriteRepository(db *sqlx.DB) payment.PaymentWriteRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetAll(ctx context.Context, limit, offset int) ([]payment.Payment, int64, error) {

	var total int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM payments`,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []payment.Payment{}, 0, nil
	}

	payments := []payment.Payment{}

	err := r.db.SelectContext(
		ctx,
		&payments,
		`SELECT * FROM payments
		 ORDER BY id
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	return payments, total, nil
}

func (r *paymentRepository) GetByBookingId(ctx context.Context, bookingID int64, limit, offset int) ([]payment.Payment, int64, error) {
	var total int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM payments 
		 WHERE booking_id = $1`,
		bookingID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []payment.Payment{}, 0, nil
	}

	properties := []payment.Payment{}
	err := r.db.SelectContext(
		ctx,
		&properties,
		`SELECT * FROM payments
		 WHERE booking_id = $1
		 ORDER BY id
		 LIMIT $2 OFFSET $3`,
		bookingID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id int64) (*payment.Payment, error) {
	var count int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT (*)
		 FROM payments
		 WHERE id = $1`,
		id,
	).Scan(&count); err != nil {
		log.Println("Error checking payment existence:", err)
		return nil, payment.ErrInternal
	}
	if count == 0 {
		return nil, payment.ErrNotFound
	}
	payments := []payment.Payment{}
	err := r.db.SelectContext(
		ctx,
		&payments,
		`SELECT * FROM payments
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		log.Println("Error fetching payment by ID:", err)
		return nil, payment.ErrInternal
	}
	payment := payments[0]
	return &payment, nil
}

func (r *paymentRepository) Create(ctx context.Context, paymentToCreate *payment.Payment) error {

	query := `
		INSERT INTO	payments (
			booking_id,
			amount,
			payment_type,
			proof_images,
			date,
			remarks,
			created_by,
			updated_by
		)
		VALUES (
			:booking_id,
			:amount,
			:payment_type,
			:proof_images,
			:date,
			:remarks,
			:created_by,
			:updated_by
		)
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, paymentToCreate)
	if err != nil {
		log.Println("Error creating payment:", err)
		return payment.ErrInternal
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&paymentToCreate.ID, &paymentToCreate.CreatedAt, &paymentToCreate.UpdatedAt)
		return nil
	}

	return payment.ErrInternal
}

func (r *paymentRepository) Update(ctx context.Context, paymentToUpdate *payment.Payment) error {
	query := `
		UPDATE payments
		SET
			amount=:amount,
			payment_type=:payment_type,
			proof_images=:proof_images,
			date=:date,
			remarks=:remarks,
			updated_by=:updated_by
		WHERE id = :id
		RETURNING updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, paymentToUpdate)
	if err != nil {
		log.Println("Error updating payment:", err)
		return payment.ErrInternal
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&paymentToUpdate.UpdatedAt)
		return nil
	}
	return nil
}
