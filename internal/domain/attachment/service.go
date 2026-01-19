package attachment

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/nevinmanoj/hostmate/internal/domain/access"
	"github.com/nevinmanoj/hostmate/internal/domain/booking"
	domaincore "github.com/nevinmanoj/hostmate/internal/domain/core"
	"github.com/nevinmanoj/hostmate/internal/domain/payment"
	"github.com/nevinmanoj/hostmate/internal/middleware"
)

type AttachmentService interface {
	RequestUploadURL(ctx context.Context, attachmentToCreate *Attachment) (int64, string, string, error)
	ConfirmUpload(ctx context.Context, attachemntID int64, success bool) (bool, string, error)
}

type attachmentService struct {
	repo          AttachmentWriteRepository
	accessService access.AccessService
	blobStorage   BlobStorage
}

func NewAttachmentService(
	repo AttachmentWriteRepository,
	accessService access.AccessService,
	blobStorage BlobStorage,
) AttachmentService {
	return &attachmentService{
		repo:          repo,
		accessService: accessService,
		blobStorage:   blobStorage,
	}
}

func (s *attachmentService) RequestUploadURL(ctx context.Context, attachmentToCreate *Attachment) (int64, string, string, error) {
	//check parent access
	userID := ctx.Value(middleware.ContextUserKey).(int64)
	err := checkParentAccess(attachmentToCreate, s.accessService, ctx)
	if err != nil {
		return -1, "", "", err
	}
	//check if the file type is valid for parent
	err = ValidateAttachmentType(attachmentToCreate.ParentType, attachmentToCreate.Type)
	if err != nil {
		return -1, "", "", err
	}

	// Generate unique blob name
	imageUUID := uuid.New().String()
	ext := filepath.Ext(attachmentToCreate.OriginalName)
	blobName := fmt.Sprintf("attachments/%d/%s%s", userID, imageUUID, ext)

	uploadURL, expiresAt, err := s.blobStorage.GenerateUploadURL(blobName)
	if err != nil {
		log.Printf("Failed to generate upload URL: %v", err)
		return -1, "", "", err
	}
	//assign rest of the variables for attachement
	attachmentToCreate.BlobName = blobName
	attachmentToCreate.Status = AttachmentStatusPending
	attachmentToCreate.CreatedBy = userID
	attachmentToCreate.UpdatedBy = userID
	attachmentToCreate.UUID = imageUUID
	err = s.repo.Create(ctx, attachmentToCreate)
	if err != nil {
		return -1, "", "", err
	}
	expiresAtStr := expiresAt.Format(time.RFC3339)
	return attachmentToCreate.ID, uploadURL, expiresAtStr, err
}

func (s *attachmentService) ConfirmUpload(ctx context.Context, attachemntID int64, success bool) (bool, string, error) {
	// //check permission
	userID := ctx.Value(middleware.ContextUserKey).(int64)

	//fetch attachement to check parent access
	attachmentFromDb, err := s.repo.GetByID(ctx, attachemntID)
	if err != nil {
		return false, "", err
	}
	err = checkParentAccess(attachmentFromDb, s.accessService, ctx)
	if err != nil {
		return false, "", err
	}

	// Verify blob exists in storage
	err = s.blobStorage.VerifyBlobExists(ctx, attachmentFromDb.BlobName)
	if err != nil {
		log.Printf("Blob not found: %v", err)
		// Update status to failed
		attachmentFromDb.Status = PropertyStatusFailed
		err = s.repo.Update(ctx, attachmentFromDb)
		if err != nil {
			return false, "", err
		}
	}

	// Update database record to "uploaded"
	attachmentFromDb.Status = AttachmentStatusUploaded
	attachmentFromDb.UpdatedBy = userID

	err = s.repo.Update(ctx, attachmentFromDb)
	if err != nil {
		return false, "", err
	}

	// Generate temporary read URL (valid for 7 days)
	readURL, err := s.blobStorage.GenerateReadURL(attachmentFromDb.BlobName)
	if err != nil {
		log.Printf("Failed to generate read URL: %v", err)
		return false, "", err
	}
	return true, readURL, nil
}

func checkParentAccess(attachment *Attachment, accessService access.AccessService, ctx context.Context) error {
	var hasAccess bool = false
	var err error
	userID, ok := ctx.Value(middleware.ContextUserKey).(int64)
	if !ok {
		return ErrInternal
	}
	switch attachment.ParentType {
	case domaincore.AttachmentParentPayment:
		hasAccess, err = accessService.CanAccessPayment(ctx, attachment.ParentId, userID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return payment.ErrUnauthorized
		}
	case domaincore.AttachmentParentBooking:
		hasAccess, err = accessService.CanAccessBooking(ctx, attachment.ParentId, userID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return booking.ErrUnauthorized
		}
	}
	return nil
}
