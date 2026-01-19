package azure

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/nevinmanoj/hostmate/internal/domain/attachment"
)

const (
	containerName = "hostmate"
	uploadExpiry  = 15 * time.Minute
	readExpiry    = 7 * 24 * time.Hour
)

func NewAzureBlobClient(connStr string) (*azblob.Client, error) {

	if connStr == "" {
		return nil, fmt.Errorf("azure connection string is empty")
	}

	client, err := azblob.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create azblob client: %w", err)
	}

	return client, nil
}

type azureBlobClient struct {
	client *azblob.Client
}

func NewBlobStorage(client *azblob.Client) attachment.BlobStorage {
	return &azureBlobClient{client: client}
}

func getBlobClient(blobName string, azureBlobClient *azblob.Client) *blob.Client {
	return azureBlobClient.ServiceClient().NewContainerClient(containerName).NewBlobClient(blobName)

}

// Generate SAS token for upload (write-only)
func (a *azureBlobClient) GenerateUploadURL(blobName string) (string, time.Time, error) {
	blobClient := getBlobClient(blobName, a.client)
	// Set expiry time
	expiresAt := time.Now().Add(uploadExpiry)

	// Create SAS token with write permission only
	sasURL, err := blobClient.GetSASURL(
		sas.BlobPermissions{Write: true, Create: true},
		expiresAt,
		nil,
		// &blob.GetSASURLOptions{
		// 	// Set content type that will be required
		// },
	)
	if err != nil {
		return "", time.Time{}, err
	}
	log.Println("SAS URL: ", sasURL)
	return sasURL, expiresAt, nil
}

func (a *azureBlobClient) GenerateReadURL(blobName string) (string, error) {
	blobClient := getBlobClient(blobName, a.client)

	// Create SAS token with read permission
	expiresAt := time.Now().Add(readExpiry)

	sasURL, err := blobClient.GetSASURL(
		sas.BlobPermissions{Read: true},
		expiresAt,
		nil,
	)
	if err != nil {
		return "", err
	}

	return sasURL, nil
}
func (a *azureBlobClient) VerifyBlobExists(ctx context.Context, blobName string) error {
	// Verify blob exists in storage
	blobClient := getBlobClient(blobName, a.client)
	_, err := blobClient.GetProperties(ctx, nil)
	return err
}
