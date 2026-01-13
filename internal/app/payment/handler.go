package payment

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	. "github.com/nevinmanoj/hostmate/api"
	errmap "github.com/nevinmanoj/hostmate/internal/app/errmap"
	payment "github.com/nevinmanoj/hostmate/internal/domain/payment"
)

type PaymentHandler struct {
	service   payment.PaymentService
	validator *validator.Validate
}

func NewPaymentHandler(s payment.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s, validator: validator.New()}
}

func (h *PaymentHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	log.Println("HandlerGetPayments::Fetching payments")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	result, totalPages, err := h.service.GetAll(r.Context(), page, pageSize)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		paymentResponses := make([]PaymentResponse, 0, len(result))
		for _, payment := range result {
			paymentResponses = append(paymentResponses, ToPaymentResponse(&payment))
		}
		resp = GetAllResponsePage[PaymentResponse]{
			StatusCode:   200,
			Message:      "Payments fetched successfully",
			TotalRecords: int64(len(result)),
			PageSize:     pageSize,
			CurrentPage:  page,
			TotalPages:   totalPages,
			Data:         paymentResponses,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) GetPaymentsWithBookingId(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerGetPaymentWithBookingId::Fetching paymenty with booking id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	bookingIdstr := chi.URLParam(r, "bookingId")
	bookingId, err := strconv.ParseInt(bookingIdstr, 10, 64)
	if err != nil {
		log.Println("handlerGetPaymentWithBookingId::Error converting property ids to int64s:", err)
		resp := errmap.GetDomainErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}

	result, totalPages, err := h.service.GetWithBookingId(r.Context(), bookingId, page, pageSize)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		paymentResponses := make([]PaymentResponse, 0, len(result))
		for _, payment := range result {
			paymentResponses = append(paymentResponses, ToPaymentResponse(&payment))
		}
		resp = GetAllResponsePage[PaymentResponse]{
			StatusCode:   200,
			Message:      "Paymnets fetched successfully for Booking ID " + bookingIdstr,
			TotalRecords: int64(len(result)),
			PageSize:     pageSize,
			CurrentPage:  page,
			TotalPages:   totalPages,
			Data:         paymentResponses,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {

	paymentIdStr := chi.URLParam(r, "paymentId")
	log.Println("HandlerGetPayment::Fetching payment with ID:", paymentIdStr)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	paymentId, err := strconv.ParseInt(paymentIdStr, 10, 64)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}
	result, err := h.service.GetById(r.Context(), paymentId)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		paymentResponse := ToPaymentResponse(result)
		resp = GetResponsePage[PaymentResponse]{
			StatusCode: 200,
			Message:    "Payment fetched successfully",
			Data:       paymentResponse,
		}
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	bookingIdStr := chi.URLParam(r, "bookingId")
	log.Println("HandlerCreatePayment::Creating payment for bookingID:", bookingIdStr)
	bookingId, err := strconv.ParseInt(bookingIdStr, 10, 64)
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var req CreatePaymentRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON body: " + err.Error(),
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
	if bookingId != req.BookingID {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "bookingID in URL and body do not match",
		})
		return
	}
	paymentToCreate := payment.Payment{
		Amount:      req.Amount,
		Date:        req.Date,
		ProofImages: req.ProofImages,
		PaymentType: req.PaymentType,
		BookingID:   req.BookingID,
		Remarks:     req.Remarks,
	}
	log.Println("Creating payment:", paymentToCreate)
	err = h.service.Create(ctx, &paymentToCreate)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)

	} else {
		paymentResponse := ToPaymentResponse(&paymentToCreate)
		resp = PostResponsePage[PaymentResponse]{
			StatusCode: http.StatusOK,
			Message:    "Payment created successfully",
			Data:       paymentResponse,
		}
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bookingIdStr := chi.URLParam(r, "bookingId")
	paymentIdStr := chi.URLParam(r, "paymentId")
	bookingId, err := strconv.ParseInt(bookingIdStr, 10, 64)
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}
	paymentId, err := strconv.ParseInt(paymentIdStr, 10, 64)

	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var req UpdatePaymentRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		log.Println("Error decoding JSON:", err)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON body",
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

	if paymentId != req.ID {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "payemnt ID in URL and body do not match or invalid",
		})
		return
	}

	if bookingId != req.BookingID {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "booking ID in URL and body do not match or invalid",
		})
		return
	}
	paymentToUpdate := payment.Payment{
		ID:          req.ID,
		Amount:      req.Amount,
		Date:        req.Date,
		ProofImages: req.ProofImages,
		PaymentType: req.PaymentType,
		BookingID:   req.BookingID,
		Remarks:     req.Remarks,
	}
	err = h.service.Update(ctx, &paymentToUpdate)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		bookingResponse := ToPaymentResponse(&paymentToUpdate)
		resp = PutResponsePage[PaymentResponse]{
			StatusCode: http.StatusOK,
			Message:    "Payment updated successfully",
			Data:       bookingResponse,
		}
		json.NewEncoder(w).Encode(resp)
	}

}
