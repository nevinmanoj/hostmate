package booking

import (
	"time"

	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
)

type CreateBookingRequest struct {
	PropertyID        int64                 `json:"property_id"`
	ManagerID         int64                 `json:"manager_id"`
	GuestPhone        string                `json:"guest_phone"`
	GuestName         string                `json:"guest_name"`
	BaseRate          float64               `json:"base_rate"`
	MaxGuestsBase     int                   `json:"max_guests_base"`
	ExtraRatePerGuest float64               `json:"extra_rate_per_guest"`
	NumGuests         int                   `json:"num_guests"`
	Status            booking.BookingStatus `json:"status"`
	CheckInDate       time.Time             `json:"check_in_date"`
	CheckOutDate      time.Time             `json:"check_out_date"`
	IdProofs          []int64               `json:"id_proofs"`
	Remarks           string                `json:"remarks"`
}
type UpdateBookingRequest struct {
	ID                int64                 `json:"id"`
	PropertyID        int64                 `json:"property_id"`
	ManagerID         int64                 `json:"manager_id"`
	GuestPhone        string                `json:"guest_phone"`
	GuestName         string                `json:"guest_name"`
	BaseRate          float64               `json:"base_rate"`
	MaxGuestsBase     int                   `json:"max_guests_base"`
	ExtraRatePerGuest float64               `json:"extra_rate_per_guest"`
	NumGuests         int                   `json:"num_guests"`
	Status            booking.BookingStatus `json:"status"`
	CheckInDate       time.Time             `json:"check_in_date"`
	CheckOutDate      time.Time             `json:"check_out_date"`
	IdProofs          []int64               `json:"id_proofs"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	CreatedBy         int64                 `json:"created_by"`
	UpdatedBy         int64                 `json:"updated_by"`
	Remarks           string                `json:"remarks"`
}
type BookingResponse struct {
	ID                int64                 `json:"id"`
	PropertyID        int64                 `json:"property_id"`
	ManagerID         int64                 `json:"manager_id"`
	GuestPhone        string                `json:"guest_phone"`
	GuestName         string                `json:"guest_name"`
	BaseRate          float64               `json:"base_rate"`
	MaxGuestsBase     int                   `json:"max_guests_base"`
	ExtraRatePerGuest float64               `json:"extra_rate_per_guest"`
	NumGuests         int                   `json:"num_guests"`
	Status            booking.BookingStatus `json:"status"`
	CheckInDate       time.Time             `json:"check_in_date"`
	CheckOutDate      time.Time             `json:"check_out_date"`
	IdProofs          []int64               `json:"id_proofs"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	CreatedBy         int64                 `json:"created_by"`
	UpdatedBy         int64                 `json:"updated_by"`
	Remarks           string                `json:"remarks"`
}

func ToBookingResponse(b *booking.Booking) BookingResponse {
	return BookingResponse{
		ID:                b.ID,
		PropertyID:        b.PropertyID,
		ManagerID:         b.ManagerID,
		GuestPhone:        b.GuestPhone,
		GuestName:         b.GuestName,
		BaseRate:          b.BaseRate,
		MaxGuestsBase:     b.MaxGuestsBase,
		ExtraRatePerGuest: b.ExtraRatePerGuest,
		NumGuests:         b.NumGuests,
		Status:            b.Status,
		CheckInDate:       b.CheckInDate,
		CheckOutDate:      b.CheckOutDate,
		IdProofs:          b.IdProofs,
		CreatedAt:         b.CreatedAt,
		UpdatedAt:         b.UpdatedAt,
		CreatedBy:         b.CreatedBy,
		UpdatedBy:         b.UpdatedBy,
		Remarks:           b.Remarks,
	}
}
