package property

import (
	"context"
)

type PropertyReadRepository interface {
	GetAll(ctx context.Context, filter PropertyFilter) ([]Property, int, error)
	GetByID(ctx context.Context, id int64) (*Property, error)
	HasManager(ctx context.Context, propertyID, userID int64) (bool, error)
}
type PropertyWriteRepository interface {
	PropertyReadRepository
	Create(ctx context.Context, property *Property) error
	Update(ctx context.Context, property *Property) error
}
