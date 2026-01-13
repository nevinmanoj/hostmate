package booking

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
)

type bookingRepository struct {
	db *sqlx.DB
}

func NewBookingReadRepository(db *sqlx.DB) booking.BookingReadRepository {
	return &bookingRepository{db: db}
}
func NewBookingWriteRepository(db *sqlx.DB) booking.BookingWriteRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) GetAll(ctx context.Context, filter booking.BookingFilter) ([]booking.Booking, int, error) {
	baseCountQuery := `SELECT COUNT(b.*)
		FROM bookings b
		JOIN properties p ON p.id = b.property_id`
	finalCountQuery, finalCountArgs, err := buildBookingQuery(baseCountQuery, filter, true)
	var total int
	if err := r.db.QueryRowContext(
		ctx,
		finalCountQuery, finalCountArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []booking.Booking{}, 0, nil
	}
	baseQuery := `SELECT b.*
		FROM bookings b
		JOIN properties p ON p.id = b.property_id`
	finalQuery, finalArgs, err := buildBookingQuery(baseQuery, filter, false)
	bookings := []booking.Booking{}
	err = r.db.SelectContext(
		ctx,
		&bookings,
		finalQuery, finalArgs...,
	)
	if err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}

func (r *bookingRepository) GetByID(ctx context.Context, id int64) (*booking.Booking, error) {
	var count int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT (*)
		 FROM bookings
		 WHERE id = $1`,
		id,
	).Scan(&count); err != nil {
		log.Println("Error checking booking existence:", err)
		return nil, booking.ErrInternal
	}
	if count == 0 {
		return nil, booking.ErrNotFound
	}
	bookings := []booking.Booking{}
	err := r.db.SelectContext(
		ctx,
		&bookings,
		`SELECT * FROM bookings
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		log.Println("Error fetching booking by ID:", err)
		return nil, booking.ErrInternal
	}
	booking := bookings[0]
	return &booking, nil
}

func (r *bookingRepository) Create(ctx context.Context, bookingToCreate *booking.Booking) error {

	query := `
		INSERT INTO bookings (
			property_id,
			manager_id,
			guest_phone,
			guest_name,
			base_rate,
			max_guests_base,
			extra_rate_per_guest,
			num_guests,
			status,
			check_in_date,
			check_out_date,
			id_proofs, 
			created_by,
			updated_by	
		)
		VALUES (
			:property_id,
			:manager_id,
			:guest_phone,
			:guest_name,
			:base_rate,
			:max_guests_base,
			:extra_rate_per_guest,
			:num_guests,
			:status,
			:check_in_date,
			:check_out_date,
			:id_proofs,
			:created_by,
			:updated_by	
		)
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, bookingToCreate)
	if err != nil {
		log.Println("Error creating booking:", err)

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// EXCLUDE constraint violation
			if pqErr.Code == "23P01" &&
				pqErr.Constraint == "no_overlapping_bookings" {
				return booking.ErrBookingConflict
			}
		}

		return booking.ErrInternal
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&bookingToCreate.ID, &bookingToCreate.CreatedAt, &bookingToCreate.UpdatedAt)
		return nil
	}

	return booking.ErrInternal
}

func (r *bookingRepository) Update(ctx context.Context, bookingToUpdate *booking.Booking) error {
	query := `
		UPDATE bookings
		SET
			property_id = :property_id,
			manager_id 	= :manager_id,
			guest_phone = :guest_phone,
			guest_name 	= :guest_name,
			base_rate 	= :base_rate,
			max_guests_base = :max_guests_base,
			extra_rate_per_guest = :extra_rate_per_guest,
			num_guests = :num_guests,
			status = :status,
			check_in_date = :check_in_date,
			check_out_date  = :check_out_date,
			id_proofs = :id_proofs, 
			updated_at = NOW(),
			updated_by = :updated_by
		WHERE id = :id
		RETURNING updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, bookingToUpdate)
	if err != nil {
		log.Println("Error updating booking:", err)
		return booking.ErrInternal
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&bookingToUpdate.UpdatedAt)
		return nil
	}
	return nil
}

func (r *bookingRepository) CheckAvailability(ctx context.Context, propertyID int64, checkInDate, checkOutDate time.Time) (bool, error) {
	available := false
	query := `
		SELECT NOT EXISTS (
    		SELECT 1
   			FROM bookings
    		WHERE property_id = 42
      		AND daterange(check_in_date, check_out_date, '[)') &&
          	daterange($2::date, $3::date, '[)')
		)`
	err := r.db.QueryRowContext(
		ctx,
		query,
		propertyID,
		checkInDate,
		checkOutDate,
	).Scan(&available)
	if err != nil {
		log.Println("Error checking booking availability:", err)
		return false, booking.ErrInternal
	}

	return available, nil
}
