package attachment

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	. "github.com/nevinmanoj/hostmate/api"
	"github.com/nevinmanoj/hostmate/internal/app/errmap"
	"github.com/nevinmanoj/hostmate/internal/domain/attachment"
)

type AttachmentHandler struct {
	service   attachment.AttachmentService
	validator *validator.Validate
}

func NewAttachmentHandler(s attachment.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: s, validator: validator.New()}
}

func (h *AttachmentHandler) RequestUploadURL(w http.ResponseWriter, r *http.Request) {
	var req ImageUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON body: " + err.Error(),
		})
		return

	}
	var resp any
	// Validate file extension
	ext := filepath.Ext(req.FileName)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		badRequestError := errmap.BadRequestError{
			Param:  "file-name",
			Reason: "Invalid file extension, must be either of ['.jpg','.jpeg','.png','.gif','.webp']",
		}
		resp = errmap.GetHttpErrorResponse(badRequestError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ctx := r.Context()

	blobName, uploadURL, expiresAt, err := h.service.RequestUploadURL(ctx, req.ParentType, req.ParentID, req.FileName)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		resp = ImageUploadResponse{
			UploadURL: uploadURL,
			BlobName:  blobName,
			ExpiresAt: expiresAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AttachmentHandler) ConfirmUpload(w http.ResponseWriter, r *http.Request) {
	var req ImageConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid JSON body: " + err.Error(),
		})
		return
	}
	ctx := r.Context()
	success, readURL, err := h.service.ConfirmUpload(ctx, req.BlobName)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}
	response := ImageConfirmResponse{
		Success:  success,
		ImageURL: readURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AttachmentHandler) ListForBooking(w http.ResponseWriter, r *http.Request) {
	h.listAttachments(w, r, attachment.AttachmentParentBooking)
}

func (h *AttachmentHandler) ListForPayment(w http.ResponseWriter, r *http.Request) {
	h.listAttachments(w, r, attachment.AttachmentParentPayment)
}
func (h *AttachmentHandler) listAttachments(w http.ResponseWriter, r *http.Request, parentType attachment.AttachmentParentType) {
	ctx := r.Context()

	parentIDStr := chi.URLParam(r, "id")
	parentID, err := strconv.ParseInt(parentIDStr, 10, 64)
	var resp any
	if err != nil {
		resp = ErrorResponse{
			StatusCode: 400,
			Message:    "Failed to fetch attachements for booking - invalid ID",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	result, err := h.service.GetAttachments(ctx, parentType, parentID)

	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		resp = GetResponsePage[[]string]{
			Message:    "Successfully genreated Read URLs for blobs of " + string(parentType),
			Data:       result,
			StatusCode: 200,
		}
	}
	json.NewEncoder(w).Encode(resp)
}
