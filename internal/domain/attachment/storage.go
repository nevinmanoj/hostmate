package attachment

import (
	"context"
	"time"
)

type BlobStorage interface {
	GenerateUploadURL(blobName string) (string, time.Time, error)
	GenerateReadURL(blobName string) (string, error)
	VerifyBlobExists(ctx context.Context, blobName string) error
	VerifyBlobSize(ctx context.Context, blobName string) error
}
