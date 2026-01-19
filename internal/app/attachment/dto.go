package attachment

import (
	attachment "github.com/nevinmanoj/hostmate/internal/domain/attachment"
	domaincore "github.com/nevinmanoj/hostmate/internal/domain/core"
)

type ImageUploadRequest struct {
	EntityType  domaincore.AttachmentParentType `json:"entity_type"`
	EntityID    int64                           `json:"entity_id"`
	FileName    string                          `json:"file_name"`
	ContentType attachment.AttchemntContentType `json:"content_type"`
	Type        attachment.AttachmentType       `json:"type"`
}

type ImageUploadResponse struct {
	UploadURL string `json:"upload_url"`
	ImageID   int64  `json:"id"`
	ExpiresAt string `json:"expiresAt"`
}

type ImageConfirmRequest struct {
	ImageID int64 `json:"id"`
	Success bool  `json:"success"`
}

type ImageConfirmResponse struct {
	Success  bool   `json:"success"`
	ImageURL string `json:"upload_url"`
}
