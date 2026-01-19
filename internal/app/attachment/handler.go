package attachment

import (
	"encoding/json"
	"net/http"
	"path/filepath"

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

	attachmentToCreate := attachment.Attachment{
		OriginalName: req.FileName,
		ContentType:  req.ContentType,
		ParentId:     req.EntityID,
		ParentType:   req.EntityType,
		Type:         req.Type,
	}
	imageID, uploadURL, expiresAt, err := h.service.RequestUploadURL(ctx, &attachmentToCreate)
	if err != nil {
		resp = errmap.GetDomainErrorResponse(err)
	} else {
		resp = ImageUploadResponse{
			UploadURL: uploadURL,
			ImageID:   imageID,
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
	success, readURL, err := h.service.ConfirmUpload(ctx, req.ImageID, req.Success)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
	response := ImageConfirmResponse{
		Success:  success,
		ImageURL: readURL,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
