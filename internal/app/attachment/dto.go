package attachment

import (
	attachment "github.com/nevinmanoj/hostmate/internal/domain/attachment"
)

type ImageUploadRequest struct {
	ParentType attachment.AttachmentParentType `json:"parent_type"`
	ParentID   int64                           `json:"parent_id"`
	FileName   string                          `json:"file_name"`
}

type ImageUploadResponse struct {
	UploadURL string `json:"upload_url"`
	BlobName  string `json:"blob_name"`
	ExpiresAt string `json:"expiresAt"`
}

type ImageConfirmRequest struct {
	BlobName string `json:"blob_name"`
}

type ImageConfirmResponse struct {
	Success  bool   `json:"success"`
	ImageURL string `json:"upload_url"`
}
