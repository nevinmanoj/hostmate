package booking

import (
	"strings"

	"github.com/jmoiron/sqlx"
	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
)

func buildBookingQuery(baseQuery string, f booking.BookingFilter, isCount bool) (string, []any, error) {
	var (
		conditions []string
		args       []any
	)

	if f.UserID != nil {
		//check if user is a manager for property,if user is nil -> admin access
		conditions = append(conditions, "? = ANY(p.managers)")
		args = append(args, *f.UserID)
	}
	if len(f.PropertyID) > 0 {
		conditions = append(conditions, "b.property_id IN (?)")
		args = append(args, f.PropertyID)
	}
	if f.ManagerID != nil {
		conditions = append(conditions, "b.manager_id = ?")
		args = append(args, f.ManagerID)
	}

	if len(f.Status) > 0 {
		conditions = append(conditions, "b.status IN (?)")
		args = append(args, f.Status)
	}

	if f.StayFrom != nil {
		conditions = append(conditions, "b.check_out_date > ?")
		args = append(args, *f.StayFrom)
	}

	if f.StayTo != nil {
		conditions = append(conditions, "b.check_in_date < ?")
		args = append(args, *f.StayTo)
	}
	if f.BookedFrom != nil {
		conditions = append(conditions, "b.created_at > ?")
		args = append(args, *f.BookedFrom)
	}

	if f.BookedTo != nil {
		conditions = append(conditions, "b.created_at < ?")
		args = append(args, *f.BookedTo)
	}

	if f.GuestPhone != nil {
		conditions = append(conditions, "b.guest_phone = ?")
		args = append(args, *f.GuestPhone)
	}

	// Apply WHERE
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Ordering (always deterministic)
	if !isCount {
		baseQuery += " ORDER BY created_at DESC"
	}

	// Pagination
	if f.Limit > 0 {
		baseQuery += " LIMIT ?"
		args = append(args, f.Limit)
	}

	if f.Offset > 0 {
		baseQuery += " OFFSET ?"
		args = append(args, f.Offset)
	}

	// Expand IN clauses
	query, finalArgs, err := sqlx.In(baseQuery, args...)
	if err != nil {
		return "", nil, err
	}

	// Rebind for postgres ($1, $2...)
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	return query, finalArgs, nil
}
