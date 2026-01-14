package property

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"

	user "github.com/nevinmanoj/hostmate/internal/domain/user"
	middleware "github.com/nevinmanoj/hostmate/internal/middleware"
)

type PropertyService interface {
	GetAll(ctx context.Context, filter PropertyFilter) ([]Property, int, error)
	GetById(ctx context.Context, id int64) (*Property, error)
	Create(ctx context.Context, property *Property) error
	Update(ctx context.Context, property *Property) error
}

type propertyService struct {
	repo     PropertyWriteRepository
	userRepo user.UserReadRepository
}

func NewPropertyService(repo PropertyWriteRepository, userRepo user.UserReadRepository) PropertyService {
	return &propertyService{repo: repo, userRepo: userRepo}
}

func (s *propertyService) GetAll(ctx context.Context, filter PropertyFilter) ([]Property, int, error) {
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	// TODO setup admin bypass
	filter.ManagerID = &userID
	data, total, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		log.Println("Error fetching properties:", err)
		return nil, 0, ErrInternal
	}
	return data, total, nil
}

func (s *propertyService) GetById(ctx context.Context, id int64) (*Property, error) {
	//check access
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	ok, err := s.repo.HasManager(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrUnauthorized
	}
	//fetch data
	data, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error fetching property with id %d: %s", id, err.Error())
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, ErrInternal
	}

	return data, nil
}

func (s *propertyService) Create(ctx context.Context, property *Property) error {
	// Validate property fields as needed, managers,images should exist
	if len(property.Managers) == 0 {
		return ErrNotValidManagers
	}
	for _, managerID := range property.Managers {
		user, err := s.userRepo.GetUserByID(ctx, managerID)
		if err != nil || user == nil {
			return ErrNotValidManagers
		}
	}
	createdBy, ok := ctx.Value(middleware.ContextUserKey).(int64)
	if !ok {
		log.Println("User not found in context")
		return ErrInternal
	}
	fmt.Printf("created by user:%d", createdBy)
	property.CreatedBy = createdBy
	property.UpdatedBy = createdBy
	fmt.Print(property)
	err := s.repo.Create(ctx, property)
	if err != nil {
		log.Println("Error creating property:", err)
		return ErrInternal
	}
	return nil
}
func (s *propertyService) Update(ctx context.Context, property *Property) error {
	//check access first
	user := ctx.Value(middleware.ContextUserKey).(int64)
	ok, err := s.repo.HasManager(ctx, property.ID, user)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUnauthorized
	}
	// Validate property fields as needed, managers,images should exist
	propertyFromDB, err := s.repo.GetByID(ctx, property.ID)
	if err != nil {
		//no such property
		return err
	}

	if len(property.Managers) == 0 {
		return ErrNotValidManagers
	}
	userInManagers := slices.Contains(property.Managers, user)
	if !userInManagers {
		log.Println("updating user must be in managers list")
		return ErrNotValidManagers
	}
	for _, managerID := range property.Managers {
		user, err := s.userRepo.GetUserByID(ctx, managerID)
		if err != nil || user == nil {
			return ErrNotValidManagers
		}
	}

	property.CreatedBy = propertyFromDB.CreatedBy
	property.CreatedAt = propertyFromDB.CreatedAt
	property.UpdatedBy = user
	err = s.repo.Update(ctx, property)
	if err != nil {
		log.Println("Error updating property:", err)
		return ErrInternal
	}
	return nil
}
