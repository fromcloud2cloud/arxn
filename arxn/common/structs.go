package common

import (
	"errors"
	"net/url"

	"github.com/google/uuid"
)

// Cloud constants.
const (
	GcpBlobName   = "gcs"
	AwsBlobName   = "s3"
	AzureBlobName = "azblob"
	Gcp           = "GCP"
	Aws           = "AWS"
	Azure         = "Azure"
)

// TransferRequest - Base Transfer Request.
type TransferRequest struct {
	TransferID       uuid.UUID
	SourceCloud      string
	DestinationCloud string
	Source           string
	Destination      string
}

// NewTransferRequest - Get a new transfer request.
func NewTransferRequest(source, destination string) (*TransferRequest, error) {
	sourceCloud, err := ParseCloudName(source)
	if err != nil {
		return &TransferRequest{}, err
	}
	destinationCloud, err := ParseCloudName(destination)
	if err != nil {
		return &TransferRequest{}, err
	}
	req := TransferRequest{
		TransferID:       uuid.New(),
		SourceCloud:      sourceCloud,
		DestinationCloud: destinationCloud,
		Source:           source,
		Destination:      destination,
	}
	return &req, nil
}

// ParseCloudName - Get cloud name from source or dest uri.
func ParseCloudName(source string) (string, error) {
	cfg := make(map[string]string)
	cfg[GcpBlobName] = Gcp
	cfg[AwsBlobName] = Aws
	cfg[AzureBlobName] = Azure

	parsedURL, err := url.Parse(source)
	if err != nil {
		return "", err
	}

	cloud, ok := cfg[parsedURL.Scheme]
	if !ok {
		err = errors.New("Invalid cloud name")
		return "", err
	}
	return cloud, nil
}
