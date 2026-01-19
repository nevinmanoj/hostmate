package attachment

import (
	"slices"
	"time"

	domaincore "github.com/nevinmanoj/hostmate/internal/domain/core"
)

type AttachmentStatus string
type AttachmentType string
type AttchemntContentType string

const (
	AttachmentStatusUploaded AttachmentStatus = "uploaded"
	AttachmentStatusPending  AttachmentStatus = "pending"
	PropertyStatusFailed     AttachmentStatus = "failed"
)
const (
	AttachmentTypePayment AttachmentType = "payment"
	AttachmentTypeIdProof AttachmentType = "id-proof"
)

const (
	AttachmentContentTypeJPEG AttchemntContentType = "image/jpeg"
	AttachmentContentTypePNG  AttchemntContentType = "image/png"
	AttachmentContentTypeGIF  AttchemntContentType = "image/gif"
	AttachmentContentTypeWEBP AttchemntContentType = "image/webp"
)

type Attachment struct {
	ID           int64                           `db:"id"`
	UUID         string                          `db:"uuid"`
	BlobName     string                          `db:"blob_name"`
	OriginalName string                          `db:"original_name"`
	ContentType  AttchemntContentType            `db:"content_type"`
	ParentId     int64                           `db:"parent_id"`
	ParentType   domaincore.AttachmentParentType `db:"parent_type"`
	Type         AttachmentType                  `db:"type"`
	Status       AttachmentStatus                `db:"status"`
	CreatedAt    time.Time                       `db:"created_at"`
	UpdatedAt    time.Time                       `db:"updated_at"`
	CreatedBy    int64                           `db:"created_by"`
	UpdatedBy    int64                           `db:"updated_by"`
}

var AllowedAttachmentTypes = map[domaincore.AttachmentParentType][]AttachmentType{
	domaincore.AttachmentParentBooking: {AttachmentTypeIdProof},
	domaincore.AttachmentParentPayment: {AttachmentTypePayment},
}

// ValidateAttachmentType checks if the attachment type is allowed for the parent type
func ValidateAttachmentType(parentType domaincore.AttachmentParentType, attachmentType AttachmentType) error {
	allowed, ok := AllowedAttachmentTypes[parentType]
	if !ok {
		return ErrInvalidAttachmentType
	}

	if slices.Contains(allowed, attachmentType) {
		return nil
	}
	return ErrInvalidAttachmentType
}
