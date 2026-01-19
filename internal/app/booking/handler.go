package booking

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	. "github.com/nevinmanoj/hostmate/api"
	errmap "github.com/nevinmanoj/hostmate/internal/app/errmap"
	booking "github.com/nevinmanoj/hostmate/internal/domain/booking"
)

type BookingHandler struct {
	service   booking.BookingService
	validator *validator.Validate
}

func NewBookingHandler(s booking.BookingService) *BookingHandler {
	return &BookingHandler{service: s, validator: validator.New()}
}

func (h *BookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	log.Println("HandlerGetBookings::Fetching bookings")

	filter, badRequestError := parseBookingFilter(r.URL.Query())
	var resp any
	if badRequestError != nil {
		resp = errmap.GetHttpErrorResponse(badRequestError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	result, total, err := h.service.GetAll(r.Context(), filter)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		bookingResponses := make([]BookingResponse, 0, len(result))
		for _, booking := range result {
			bookingResponses = append(bookingResponses, ToBookingResponse(&booking))
		}
		resp = GetAllResponsePage[BookingResponse]{
			StatusCode:   200,
			Message:      "Bookings fetched successfully",
			TotalRecords: total,
			Limit:        filter.Limit,
			Offset:       filter.Offset,
			Data:         bookingResponses,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "bookingId")
	log.Println("HandlerGetBooking::Fetching booking with ID:", idStr)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp = ErrorResponse{
			StatusCode: 500,
			Message:    "Failed to fetch booking - invalid ID",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	result, err := h.service.GetById(r.Context(), id)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		bookingResponse := ToBookingResponse(result)
		resp = GetResponsePage[BookingResponse]{
			StatusCode: 200,
			Message:    "Booking fetched successfully",
			Data:       bookingResponse,
		}
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateBookingRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid JSON body: " + err.Error(),
		})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	bookingToCreate := booking.Booking{
		PropertyID:        req.PropertyID,
		ManagerID:         req.ManagerID,
		GuestPhone:        req.GuestPhone,
		GuestName:         req.GuestName,
		BaseRate:          req.BaseRate,
		MaxGuestsBase:     req.MaxGuestsBase,
		ExtraRatePerGuest: req.ExtraRatePerGuest,
		NumGuests:         req.NumGuests,
		Status:            req.Status,
		CheckInDate:       NormalizeDate(req.CheckInDate),
		CheckOutDate:      NormalizeDate(req.CheckOutDate),
		Remarks:           req.Remarks,
	}
	fmt.Println("Creating booking:", bookingToCreate)
	err := h.service.Create(ctx, &bookingToCreate)
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)

	} else {
		bookingResponse := ToBookingResponse(&bookingToCreate)
		resp = PostResponsePage[BookingResponse]{
			StatusCode: http.StatusOK,
			Message:    "Booking created successfully",
			Data:       bookingResponse,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req UpdateBookingRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		fmt.Println("Error decoding JSON:", err)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid JSON body",
		})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	idStr := chi.URLParam(r, "bookingId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id != req.ID {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "ID in URL and body do not match or invalid",
		})
		return
	}
	bookingToUpdate := booking.Booking{
		ID:                req.ID,
		PropertyID:        req.PropertyID,
		ManagerID:         req.ManagerID,
		GuestPhone:        req.GuestPhone,
		GuestName:         req.GuestName,
		BaseRate:          req.BaseRate,
		MaxGuestsBase:     req.MaxGuestsBase,
		ExtraRatePerGuest: req.ExtraRatePerGuest,
		NumGuests:         req.NumGuests,
		Status:            req.Status,
		CheckInDate:       req.CheckInDate,
		CheckOutDate:      req.CheckOutDate,
		Remarks:           req.Remarks,
	}
	err = h.service.Update(ctx, &bookingToUpdate)
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		bookingResponse := ToBookingResponse(&bookingToUpdate)
		resp = PutResponsePage[BookingResponse]{
			StatusCode: http.StatusOK,
			Message:    "Property updated successfully",
			Data:       bookingResponse,
		}
		json.NewEncoder(w).Encode(resp)
	}

}

func (h *BookingHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "propertyId")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	startDateTime, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid start_date format, expected YYYY-MM-DD",
		})
		return
	}
	endDateTime, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid end_date format, expected YYYY-MM-DD",
		})
		return
	}
	log.Println("HandlerCheckAvailability::Checking availability for property ID:", idStr)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}
	available, err := h.service.CheckAvailability(r.Context(), id, startDateTime, endDateTime)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		resp = GetResponsePage[bool]{
			StatusCode: 200,
			Message:    "Availability checked successfully",
			Data:       available,
		}
	}
}

func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
