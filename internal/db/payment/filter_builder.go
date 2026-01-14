package payment

import (
	"strings"

	"github.com/jmoiron/sqlx"
	payment "github.com/nevinmanoj/hostmate/internal/domain/payment"
)

func buildPaymentQuery(baseQuery string, f payment.PaymentFilter, isCount bool) (string, []any, error) {
	var (
		conditions []string
		args       []any
	)

	if f.UserID != nil {
		conditions = append(conditions, "? = ANY(pr.managers)")
		args = append(args, *f.UserID)
	}

	if len(f.PaymentType) > 0 {
		conditions = append(conditions, "p.payment_type IN (?)")
		args = append(args, f.PaymentType)
	}
	if f.FromDate != nil {
		conditions = append(conditions, "p.created_at > ?")
		args = append(args, *f.FromDate)
	}

	if f.ToDate != nil {
		conditions = append(conditions, "p.created_at < ?")
		args = append(args, *f.ToDate)
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
