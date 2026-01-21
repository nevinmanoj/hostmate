package attachment

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nevinmanoj/hostmate/internal/domain/access"
	"github.com/nevinmanoj/hostmate/internal/domain/booking"
	"github.com/nevinmanoj/hostmate/internal/domain/payment"
	"github.com/nevinmanoj/hostmate/internal/middleware"
)

type AttachmentService interface {
	RequestUploadURL(ctx context.Context, parentType AttachmentParentType, parentID int64, fileName string) (string, string, string, error)
	ConfirmUpload(ctx context.Context, blobName string) (bool, string, error)
	GetAttachments(ctx context.Context, parentType AttachmentParentType, parentID int64) ([]string, error)
}

type attachmentService struct {
	accessService  access.AccessService
	blobStorage    BlobStorage
	paymentService payment.PaymentService
	bookingService booking.BookingService
}

func NewAttachmentService(
	accessService access.AccessService,
	blobStorage BlobStorage,
	paymentService payment.PaymentService,
	bookingService booking.BookingService,

) AttachmentService {
	return &attachmentService{
		accessService:  accessService,
		blobStorage:    blobStorage,
		paymentService: paymentService,
		bookingService: bookingService,
	}
}

func (s *attachmentService) RequestUploadURL(ctx context.Context, parentType AttachmentParentType, parentID int64, fileName string) (string, string, string, error) {
	//check parent access
	err := checkParentAccess(parentType, parentID, s.accessService, ctx)
	if err != nil {
		return "", "", "", err
	}

	// Generate unique blob name
	imageUUID := uuid.New().String()
	ext := strings.ToLower(filepath.Ext(fileName))
	blobName := fmt.Sprintf("%s/%d/%s%s", string(parentType), parentID, imageUUID, ext)

	uploadURL, expiresAt, err := s.blobStorage.GenerateUploadURL(blobName)
	if err != nil {
		log.Printf("Failed to generate upload URL: %v", err)
		return "", "", "", err
	}

	expiresAtStr := expiresAt.Format(time.RFC3339)
	return blobName, uploadURL, expiresAtStr, err
}
func (s *attachmentService) ConfirmUpload(ctx context.Context, blobName string) (bool, string, error) {
	//extract parentType and parentID from blobname
	parentType, parentID, err := ParseBlobName(blobName)
	if err != nil {
		return false, "", err
	}
	// Verify blob exists in storage and the size is valid
	err = s.blobStorage.VerifyBlobSize(ctx, blobName)
	if err != nil {
		log.Printf("invalid blob: %v", err)
		return false, "", err
	}

	//update new blobName in parent blobs array
	switch parentType {
	case AttachmentParentBooking:
		err = s.bookingService.ConfirmBlobsUpload(ctx, parentID, blobName)
	case AttachmentParentPayment:
		err = s.paymentService.ConfirmBlobsUpload(ctx, parentID, blobName)
	}
	if err != nil {
		log.Printf("Failed to update parent blobs: %v", err)
		return false, "", err
	}
	// Generate temporary read URL (valid for 7 days)
	readURL, err := s.blobStorage.GenerateReadURL(blobName)
	if err != nil {
		log.Printf("Failed to generate read URL: %v", err)
		return false, "", err
	}
	return true, readURL, nil
}
func (s *attachmentService) GetAttachments(ctx context.Context, parentType AttachmentParentType, parentID int64) ([]string, error) {

	//fetch blobs, access is checked in parent service
	blobs := []string{}
	var err error
	switch parentType {
	case AttachmentParentBooking:
		blobs, err = s.bookingService.GetBlobs(ctx, parentID)
	case AttachmentParentPayment:
		blobs, err = s.paymentService.GetBlobs(ctx, parentID)
	}
	if err != nil {
		return nil, err
	}
	blobURLs := []string{}
	//process blobs into read urls
	for _, blob := range blobs {
		url, err := s.blobStorage.GenerateReadURL(blob)
		if err == nil {
			blobURLs = append(blobURLs, url)
		}
	}
	return blobURLs, nil
}

// helpers
func checkParentAccess(parentType AttachmentParentType, parentID int64, accessService access.AccessService, ctx context.Context) error {
	var hasAccess bool = false
	var err error
	userID, ok := ctx.Value(middleware.ContextUserKey).(int64)
	if !ok {
		return ErrInternal
	}
	switch parentType {
	case AttachmentParentPayment:
		hasAccess, err = accessService.CanAccessPayment(ctx, parentID, userID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return payment.ErrUnauthorized
		}
	case AttachmentParentBooking:
		hasAccess, err = accessService.CanAccessBooking(ctx, parentID, userID)
		if err != nil {
			return err
		}
		if !hasAccess {
			return booking.ErrUnauthorized
		}
	}
	return nil
}

func ParseBlobName(blobName string) (parentType AttachmentParentType, parentID int64, err error) {
	parts := strings.Split(blobName, "/")
	log.Println("parts:", parts)
	if len(parts) != 3 {
		return "", 0, fmt.Errorf("invalid blobName format: %s", blobName)
	}

	// Parse parentType
	parentType = AttachmentParentType(parts[0])
	switch parentType {
	case AttachmentParentBooking, AttachmentParentPayment:
		// valid
	default:
		return "", 0, ErrInvalidAttachmentParentType
	}

	// Parse parentID
	parentID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil || parentID < 0 {
		return "", 0, fmt.Errorf("invalid parent ID in blobName: %s", parts[1])
	}

	// Validate filename part: "<uuid><ext>"
	filePart := parts[2]
	ext := filepath.Ext(filePart)
	if ext == "" {
		return "", 0, fmt.Errorf("missing file extension in blobName: %s", blobName)
	}

	name := strings.TrimSuffix(filePart, ext)
	if name == "" {
		return "", 0, fmt.Errorf("invalid file name in blobName: %s", blobName)
	}

	// Optional: strict UUID validation
	if _, uuidErr := uuid.Parse(name); uuidErr != nil {
		return "", 0, fmt.Errorf("invalid UUID in blobName: %s", name)
	}

	return parentType, parentID, nil
}
