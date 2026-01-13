package property

import (
	"strings"

	"github.com/jmoiron/sqlx"
	property "github.com/nevinmanoj/hostmate/internal/domain/property"
)

func buildPropertyQuery(baseQuery string, f property.PropertyFilter, isCount bool) (string, []any, error) {
	var (
		conditions []string
		args       []any
	)

	// if f.UserID != nil {
	// 	conditions = append(conditions, "user_id = ?")
	// 	args = append(args, *f.UserID)
	// }

	if len(f.Type) > 0 {
		conditions = append(conditions, "type IN (?)")
		args = append(args, f.Type)
	}

	if f.ManagerID != nil {
		conditions = append(conditions, "? = ANY(managers)")
		args = append(args, f.ManagerID)
	}

	if f.Active != nil {
		conditions = append(conditions, "active = (?)")
		args = append(args, *f.Active)
	}

	// Apply WHERE
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Ordering (always deterministic)
	if !isCount {
		baseQuery += " ORDER BY created_at DESC"
		// Pagination
		if f.Limit > 0 {
			baseQuery += " LIMIT ?"
			args = append(args, f.Limit)
		}

		if f.Offset > 0 {
			baseQuery += " OFFSET ?"
			args = append(args, f.Offset)
		}
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
