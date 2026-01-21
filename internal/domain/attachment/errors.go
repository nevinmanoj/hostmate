package attachment

import "errors"

var (
	ErrInternal                    = errors.New("Internal error")
	ErrInvalidAttachmentParentType = errors.New("invalid attachment parent type")
)
