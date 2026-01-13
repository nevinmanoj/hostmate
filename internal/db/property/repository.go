package property

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	property "github.com/nevinmanoj/hostmate/internal/domain/property"
)

type propertyRepository struct {
	db *sqlx.DB
}

func NewPropertyWriteRepository(db *sqlx.DB) property.PropertyWriteRepository {
	return &propertyRepository{db: db}
}
func NewPropertyReadRepository(db *sqlx.DB) property.PropertyReadRepository {
	return &propertyRepository{db: db}
}
func (r *propertyRepository) GetAll(ctx context.Context, filter property.PropertyFilter) ([]property.Property, int, error) {

	baseCountQuery := `SELECT COUNT(*) FROM properties`
	finalCountQuery, finalCountArgs, err := buildPropertyQuery(baseCountQuery, filter, true)
	if err != nil {
		log.Println("Error during building properties query:", err.Error())
		return nil, 0, err
	}

	var total int
	if err := r.db.QueryRowContext(
		ctx,
		finalCountQuery, finalCountArgs...,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []property.Property{}, 0, nil
	}
	baseQuery := `SELECT * FROM properties`
	finalQuery, finalArgs, err := buildPropertyQuery(baseQuery, filter, false)
	properties := []property.Property{}
	err = r.db.SelectContext(
		ctx,
		&properties,
		finalQuery, finalArgs...,
	)
	if err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}
func (r *propertyRepository) GetByManagerId(ctx context.Context, managerID int64, limit, offset int) ([]property.Property, int64, error) {

	var total int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM properties 
		 WHERE $1 = ANY(managers)`,
		managerID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []property.Property{}, 0, nil
	}

	properties := []property.Property{}
	err := r.db.SelectContext(
		ctx,
		&properties,
		`SELECT * FROM properties
		 WHERE $1 = ANY(managers)
		 ORDER BY id
		 LIMIT $2 OFFSET $3`,
		managerID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

func (r *propertyRepository) GetByID(ctx context.Context, id int64) (*property.Property, error) {
	var count int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT (*) 
		 FROM properties 
		 WHERE id = $1`,
		id,
	).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, property.ErrNotFound
	}
	properties := []property.Property{}
	err := r.db.SelectContext(
		ctx,
		&properties,
		`SELECT * FROM properties
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}
	property := properties[0]
	return &property, nil
}
func (r *propertyRepository) Create(ctx context.Context, propertyToCreate *property.Property) error {

	query := `
		INSERT INTO properties (
			name,
			address,
			type,
			base_rate,
			max_guests_base,
			extra_rate_per_guest,
			managers,
			photos,
			active,
			created_by,
			updated_by
		)
		VALUES (
			:name,
			:address,
			:type,
			:base_rate,
			:max_guests_base,
			:extra_rate_per_guest,
			:managers,
			:photos,
			:active,
			:created_by,
			:updated_by
		)
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, propertyToCreate)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&propertyToCreate.ID, &propertyToCreate.CreatedAt, &propertyToCreate.UpdatedAt)
		return nil
	}

	return sql.ErrNoRows
}
func (r *propertyRepository) Update(ctx context.Context, propertyToUpdate *property.Property) error {
	query := `
		UPDATE properties
		SET
			name = :name,
			address = :address,
			type = :type,
			base_rate = :base_rate,
			max_guests_base = :max_guests_base,
			extra_rate_per_guest = :extra_rate_per_guest,
			managers = :managers,
			photos = :photos,
			active = :active,
			updated_at = NOW(),
			updated_by = :updated_by
		WHERE id = :id
		RETURNING updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, propertyToUpdate)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&propertyToUpdate.UpdatedAt)
		return nil
	}
	return nil
}
