package booking

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	errMap "github.com/nevinmanoj/hostmate/internal/app/errmap"
	httputil "github.com/nevinmanoj/hostmate/internal/app/httputil"
	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
)

func parseBookingFilter(q url.Values) (booking.BookingFilter, *errMap.BadRequestError) {
	var f booking.BookingFilter

	if v := q.Get("property_id"); v != "" {
		propertyId, err := httputil.ParseInt64Slice(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "property_id",
				Reason: err.Error(),
			}
		}
		f.PropertyID = propertyId
	}

	if v := q.Get("status"); v != "" {
		status, err := parseBookingStatusSlice(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "status",
				Reason: err.Error(),
			}
		}
		f.Status = status
	}
	if v := q.Get("stay_from"); v != "" {
		stayFrom, err := httputil.ParseDatePtr(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "stay_from",
				Reason: "invalid date format, expected YYYY-MM-DD",
			}
		}
		f.StayFrom = stayFrom
	}

	if v := q.Get("stay_to"); v != "" {
		stayTo, err := httputil.ParseDatePtr(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "stay_to",
				Reason: "invalid date format, expected YYYY-MM-DD",
			}
		}
		f.StayTo = stayTo
	}
	if v := q.Get("booked_from"); v != "" {
		stayFrom, err := httputil.ParseDatePtr(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "booked_from",
				Reason: "invalid date format, expected YYYY-MM-DD",
			}
		}
		f.StayFrom = stayFrom
	}

	if v := q.Get("booked_to"); v != "" {
		stayTo, err := httputil.ParseDatePtr(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "booked_to",
				Reason: "invalid date format, expected YYYY-MM-DD",
			}
		}
		f.StayTo = stayTo
	}

	// Pagination defaults
	f.Limit = 100
	f.Offset = 0

	if v := q.Get("limit"); v != "" {
		limit, err := strconv.Atoi(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "limit",
				Reason: err.Error(),
			}
		} else if limit > 0 && limit < 100 {
			f.Limit = limit
		}
	}

	if v := q.Get("offset"); v != "" {
		offset, err := strconv.Atoi(v)
		if err != nil {
			return f, &errMap.BadRequestError{
				Param:  "offset",
				Reason: err.Error(),
			}
		} else if offset > 0 {
			f.Offset = offset
		}
	}

	return f, nil
}

func ParseBookingStatus(v string) (booking.BookingStatus, error) {
	switch strings.ToLower(v) {

	case "booked":
		return booking.BookingBooked, nil
	case "checkedIn":
		return booking.BookingCheckedIn, nil
	case "checkedOut":
		return booking.BookingCheckedOut, nil
	case "cancelled":
		return booking.BookingCancelled, nil
	default:
		return "", fmt.Errorf("Invalid type, must be ['booked','checkedIn','checkedOut','cancelled'] ")
	}
}

func parseBookingStatusSlice(v string) ([]booking.BookingStatus, error) {
	parts := strings.Split(v, ",")
	out := make([]booking.BookingStatus, 0, len(parts))

	for _, p := range parts {
		t, err := ParseBookingStatus(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}

	return out, nil
}
