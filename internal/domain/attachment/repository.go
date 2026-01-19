package attachment

import "context"

type AttachmentReadRepository interface {
	GetByID(ctx context.Context, id int64) (*Attachment, error)
}

type AttachmentWriteRepository interface {
	AttachmentReadRepository
	Create(ctx context.Context, attachmentToCreate *Attachment) error
	Update(ctx context.Context, attachmentToUpdate *Attachment) error
}
