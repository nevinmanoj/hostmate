package access

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nevinmanoj/hostmate/internal/domain/access"
)

type accessRepository struct {
	db *sqlx.DB
}

func NewAccessRepository(db *sqlx.DB) access.AccessRepository {
	return &accessRepository{db: db}
}

func (r *accessRepository) HasManagerByPropertyID(ctx context.Context, propertyID, userID int64) (bool, error) {

	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM properties pr
			WHERE id = $1
			  AND $2 = ANY(pr.managers)
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, q, propertyID, userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (r *accessRepository) HasManagerByBookingID(ctx context.Context, bookingID, userID int64) (bool, error) {

	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM bookings b
			JOIN properties pr ON pr.id = b.property_id
			WHERE b.id = $1
			  AND $2 = ANY(pr.managers)
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, q, bookingID, userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
func (r *accessRepository) HasManagerByPaymentID(ctx context.Context, paymentID, userID int64) (bool, error) {

	const q = `
		SELECT EXISTS (
			SELECT 1
			FROM payments p
			JOIN bookings b ON b.id = p.booking_id
			JOIN properties pr ON pr.id = b.property_id
			WHERE p.id = $1
			  AND $2 = ANY(pr.managers)
		)
	`

	var exists bool
	err := r.db.GetContext(ctx, &exists, q, paymentID, userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
