package property

import (
	property "github.com/nevinmanoj/hostmate/internal/domain/property"
)

type CreatePropertyRequest struct {
	Name              string   `json:"name" validate:"required,min=2"`
	Address           string   `json:"address" validate:"required"`
	Type              string   `json:"type" validate:"required,oneof=apartment villa hotel"`
	BaseRate          float64  `json:"base_rate" validate:"gt=0"`
	MaxGuestsBase     int      `json:"max_guests_base" validate:"gt=0"`
	ExtraRatePerGuest float64  `json:"extra_rate_per_guest" validate:"gte=0"`
	Managers          []int64  `json:"managers" validate:"required,min=1,dive,gt=0"`
	Photos            []string `json:"photos" validate:"omitempty,dive,url"`
	Active            *bool    `json:"active,omitempty"`
}
type UpdatePropertyRequest struct {
	ID                int64    `json:"id" validate:"required,gt=0"`
	Name              string   `json:"name" validate:"required,min=2"`
	Address           string   `json:"address" validate:"required"`
	Type              string   `json:"type" validate:"required,oneof=apartment villa room"`
	BaseRate          float64  `json:"base_rate" validate:"gt=0"`
	MaxGuestsBase     int      `json:"max_guests_base" validate:"gt=0"`
	ExtraRatePerGuest float64  `json:"extra_rate_per_guest" validate:"gte=0"`
	Managers          []int64  `json:"managers" validate:"required,min=1,dive,gt=0"`
	Photos            []string `json:"photos" validate:"omitempty"`
	Active            *bool    `json:"active,omitempty"`
}

type PropertyResponse struct {
	ID                int64    `json:"id"`
	Name              string   `json:"name"`
	Address           string   `json:"address" `
	Type              string   `json:"type"`
	BaseRate          float64  `json:"base_rate"`
	MaxGuestsBase     int      `json:"max_guests_base"`
	ExtraRatePerGuest float64  `json:"extra_rate_per_guest"`
	Managers          []int64  `json:"managers" validate:"required,min=1,dive,gt=0"`
	Photos            []string `json:"photos" validate:"omitempty"`
	Active            bool     `json:"active"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
	UpdatedBy         int64    `json:"updated_by"`
	CreatedBy         int64    `json:"created_by"`
}

func ToPropertyResponse(p *property.Property) PropertyResponse {
	return PropertyResponse{
		ID:                p.ID,
		Name:              p.Name,
		Address:           p.Address,
		Type:              string(p.Type),
		BaseRate:          p.BaseRate,
		MaxGuestsBase:     p.MaxGuestsBase,
		ExtraRatePerGuest: p.ExtraRatePerGuest,
		Managers:          p.Managers,
		Photos:            p.Photos,
		Active:            p.Active,
		CreatedAt:         p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:         p.CreatedBy,
		UpdatedBy:         p.UpdatedBy,
	}
}
