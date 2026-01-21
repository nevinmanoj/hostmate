package booking

import (
	"time"
)

type BookingStatus string

const (
	BookingBooked     BookingStatus = "booked"
	BookingCheckedIn  BookingStatus = "checkedIn"
	BookingCheckedOut BookingStatus = "checkedOut"
	BookingCancelled  BookingStatus = "cancelled"
)

type Booking struct {
	ID                int64         `db:"id"`
	PropertyID        int64         `db:"property_id"`
	ManagerID         int64         `db:"manager_id"`
	GuestPhone        string        `db:"guest_phone"`
	GuestName         string        `db:"guest_name"`
	BaseRate          float64       `db:"base_rate"`
	MaxGuestsBase     int           `db:"max_guests_base"`
	ExtraRatePerGuest float64       `db:"extra_rate_per_guest"`
	NumGuests         int           `db:"num_guests"`
	Status            BookingStatus `db:"status"`
	blobs             []string      `db:"blobs"`
	CheckInDate       time.Time     `db:"check_in_date"`
	CheckOutDate      time.Time     `db:"check_out_date"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
	CreatedBy         int64         `db:"created_by"`
	UpdatedBy         int64         `db:"updated_by"`
	Remarks           string        `db:"remarks"`
}
