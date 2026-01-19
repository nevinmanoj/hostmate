package attachment

import "errors"

var (
	ErrNotFound                    = errors.New("Attachment not found")
	ErrInternal                    = errors.New("Internal error")
	ErrUnauthorized                = errors.New("unauthorized")
	ErrInvalidAttachmentType       = errors.New("invalid attachment type for entity")
	ErrInvalidAttachmentParentType = errors.New("invalid attachment entity type")
)
