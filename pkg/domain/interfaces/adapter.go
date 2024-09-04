package interfaces

import (
	"context"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AzureBlobStorage interface {
	NewReader(ctx context.Context, storageAccountName, containerName, blobName string) (io.ReadCloser, error)
}

type GoogleCloudStorage interface {
	NewWriter(ctx context.Context, bucketName, objectName string) (io.WriteCloser, error)
}
