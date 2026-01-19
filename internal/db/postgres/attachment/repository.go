package attachment

import (
	"context"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/nevinmanoj/hostmate/internal/domain/attachment"
)

type attachmentRepository struct {
	db *sqlx.DB
}

func NewAttachmentWriteRepository(db *sqlx.DB) attachment.AttachmentWriteRepository {
	return &attachmentRepository{db: db}
}
func NewAttachmentReadRepository(db *sqlx.DB) attachment.AttachmentReadRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(ctx context.Context, attachmentToCreate *attachment.Attachment) error {

	query := `
		INSERT INTO attachments (
			uuid,
			blob_name,
			original_name,
			content_type,
			parent_id,
			parent_type,
			type,
			status,
			created_by,
			updated_by
		)
		VALUES (
			:uuid,
			:blob_name,
			:original_name,
			:content_type,
			:parent_id,
			:parent_type,
			:type,
			:status,
			:created_by,
			:updated_by
		)
		RETURNING id, created_at, updated_at
	`
	rows, err := r.db.NamedQueryContext(ctx, query, attachmentToCreate)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&attachmentToCreate.ID, &attachmentToCreate.CreatedAt, &attachmentToCreate.UpdatedAt)
		return nil
	}

	return sql.ErrNoRows
}

func (r *attachmentRepository) Update(ctx context.Context, attachmentToUpdate *attachment.Attachment) error {
	query := `
		UPDATE attachments
		SET
			uuid=:uuid,
			blob_name = :blob_name,
			original_name = :original_name,
			content_type = :content_type,
			parent_id = :parent_id,
			parent_type = :parent_type,
			type = :type,
			status = :status,
			updated_at = NOW(),
			updated_by = :updated_by
		
		WHERE id = :id
		RETURNING updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, attachmentToUpdate)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&attachmentToUpdate.UpdatedAt)
		return nil
	}
	return nil
}

func (r *attachmentRepository) GetByID(ctx context.Context, id int64) (*attachment.Attachment, error) {
	var count int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT (*)
		 FROM attachments
		 WHERE id = $1`,
		id,
	).Scan(&count); err != nil {
		log.Println("Error checking attachment's existence:", err)
		return nil, attachment.ErrInternal
	}
	if count == 0 {
		return nil, attachment.ErrNotFound
	}
	attachments := []attachment.Attachment{}
	err := r.db.SelectContext(
		ctx,
		&attachments,
		`SELECT * FROM attachments
		 WHERE id = $1`,
		id,
	)
	if err != nil {
		log.Println("Error fetching attachment by ID:", err)
		return nil, attachment.ErrInternal
	}
	attachment := attachments[0]
	return &attachment, nil
}
