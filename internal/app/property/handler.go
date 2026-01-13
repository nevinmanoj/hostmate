package property

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"

	. "github.com/nevinmanoj/hostmate/api"
	errmap "github.com/nevinmanoj/hostmate/internal/app/errmap"
	property "github.com/nevinmanoj/hostmate/internal/domain/property"
)

type PropertyHandler struct {
	service   property.PropertyService
	validator *validator.Validate
}

func NewPropertyHandler(s property.PropertyService) *PropertyHandler {
	return &PropertyHandler{service: s, validator: validator.New()}
}

func (h *PropertyHandler) GetProperties(w http.ResponseWriter, r *http.Request) {
	log.Println("HandlerGetProperties::Fetching properties")
	filter, badRequestError := parsePropertyFilter(r.URL.Query())
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
		propertyResponses := make([]PropertyResponse, 0, len(result))
		for _, property := range result {
			propertyResponses = append(propertyResponses, ToPropertyResponse(&property))
		}
		resp = GetAllnewResponsePage[PropertyResponse]{
			StatusCode:   200,
			Message:      "Properties fetched successfully",
			TotalRecords: total,
			Limit:        filter.Limit,
			Offset:       filter.Offset,
			Data:         propertyResponses,
		}
	}
	json.NewEncoder(w).Encode(resp)
}
func (h *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "propertyId")
	log.Println("HandlerGetProperty::Fetching property with ID:", idStr)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}
	result, err := h.service.GetById(r.Context(), id)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		propertyResponse := ToPropertyResponse(result)
		resp = GetResponsePage[PropertyResponse]{
			StatusCode: 200,
			Message:    "Property fetched successfully",
			Data:       propertyResponse,
		}
	}

	json.NewEncoder(w).Encode(resp)
}
func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreatePropertyRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(&req); err != nil {
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
	propertyToCreate := property.Property{
		Name:              req.Name,
		Address:           req.Address,
		Type:              property.PropertyType(req.Type),
		BaseRate:          req.BaseRate,
		MaxGuestsBase:     req.MaxGuestsBase,
		ExtraRatePerGuest: req.ExtraRatePerGuest,
		Managers:          req.Managers,
		Photos:            req.Photos,
		Active:            true,
	}
	if req.Active != nil {
		propertyToCreate.Active = *req.Active
	}
	fmt.Println("Creating property:", propertyToCreate)
	err := h.service.Create(ctx, &propertyToCreate)
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	}
	propertyResponse := ToPropertyResponse(&propertyToCreate)
	resp = PostResponsePage[PropertyResponse]{
		StatusCode: http.StatusOK,
		Message:    "Property created successfully",
		Data:       propertyResponse,
	}
	json.NewEncoder(w).Encode(resp)
}
func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req UpdatePropertyRequest
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
	idStr := chi.URLParam(r, "propertyId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id != req.ID {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "ID in URL and body do not match or invalid",
		})
		return
	}
	propertyToUpdate := property.Property{
		ID:                req.ID,
		Name:              req.Name,
		Address:           req.Address,
		Type:              property.PropertyType(req.Type),
		BaseRate:          req.BaseRate,
		MaxGuestsBase:     req.MaxGuestsBase,
		ExtraRatePerGuest: req.ExtraRatePerGuest,
		Managers:          req.Managers,
		Photos:            req.Photos,
		Active:            *req.Active,
	}
	err = h.service.Update(ctx, &propertyToUpdate)
	var resp any
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		propertyResponse := ToPropertyResponse(&propertyToUpdate)
		resp = PutResponsePage[PropertyResponse]{
			StatusCode: 200,
			Message:    "Property updated successfully",
			Data:       propertyResponse,
		}
	}
	json.NewEncoder(w).Encode(resp)
}
