package property

import (
	"context"
)

type PropertyReadRepository interface {
	GetAll(ctx context.Context, filter PropertyFilter) ([]Property, int, error)
	GetByManagerId(ctx context.Context, managerID int64, limit, offset int) ([]Property, int64, error)
	GetByID(ctx context.Context, id int64) (*Property, error)
}
type PropertyWriteRepository interface {
	PropertyReadRepository
	Create(ctx context.Context, property *Property) error
	Update(ctx context.Context, property *Property) error
}
