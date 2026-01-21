package errmap

import (
	"fmt"

	. "github.com/nevinmanoj/hostmate/api"
	"github.com/nevinmanoj/hostmate/internal/domain/attachment"
	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
	"github.com/nevinmanoj/hostmate/internal/domain/payment"
	property "github.com/nevinmanoj/hostmate/internal/domain/property"
	user "github.com/nevinmanoj/hostmate/internal/domain/user"
)

func GetDomainErrorResponse(err error) ErrorResponse {
	switch err {
	//user errrors
	case user.ErrUnauthorized:
		return ErrorResponse{
			StatusCode: 403,
			Message:    "Unauthorized to access user",
		}
	case user.ErrNotFound:
		return ErrorResponse{
			StatusCode: 404,
			Message:    "User not found",
		}
	case user.ErrAlreadyExists:
		return ErrorResponse{
			StatusCode: 400,
			Message:    "User already exists",
		}
	//property errors
	case property.ErrUnauthorized:
		return ErrorResponse{
			StatusCode: 403,
			Message:    "Unauthorized to access properties",
		}
	case property.ErrNotFound:
		return ErrorResponse{
			StatusCode: 404,
			Message:    "Property not found",
		}
	case property.ErrNotValidManagers:
		return ErrorResponse{
			StatusCode: 400,
			Message:    "Managers are not valid",
		}
	//booking errors
	case booking.ErrUnauthorized:
		return ErrorResponse{
			StatusCode: 403,
			Message:    "Unauthorized to access this booking",
		}
	case booking.ErrNotFound:
		return ErrorResponse{
			StatusCode: 404,
			Message:    "Booking not found",
		}
	case booking.ErrInvalidDateRange:
		return ErrorResponse{
			StatusCode: 400,
			Message:    "The provided date range is invalid",
		}
	case booking.ErrBookingConflict:
		return ErrorResponse{
			StatusCode: 409,
			Message:    "The booking dates conflict with an existing booking",
		}
	//payments
	case payment.ErrUnauthorized:
		return ErrorResponse{
			StatusCode: 403,
			Message:    "Unauthorized to access this payment",
		}
	case payment.ErrNotValidBookingId:
		return ErrorResponse{
			StatusCode: 400,
			Message:    "Invalid booking_id",
		}
	case payment.ErrNotFound:
		return ErrorResponse{
			StatusCode: 404,
			Message:    "Payment not found",
		}
	//attachments
	case attachment.ErrInvalidAttachmentParentType:
		return ErrorResponse{
			StatusCode: 400,
			Message:    "Invalid value for attachment parent type",
		}
	default:
		return ErrorResponse{
			StatusCode: 500,
			Message:    "Error: " + fmt.Sprint(err),
		}
	}

}
