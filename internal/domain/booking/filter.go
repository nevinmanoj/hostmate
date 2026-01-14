package booking

import "time"

type BookingFilter struct {
	Id         *int64
	UserID     *int64
	PropertyID []int64
	Status     []BookingStatus
	ManagerID  *int64
	BookedFrom *time.Time
	BookedTo   *time.Time
	StayFrom   *time.Time
	StayTo     *time.Time
	GuestPhone *string
	Limit      int
	Offset     int
}
