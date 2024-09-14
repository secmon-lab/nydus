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
	NewReader(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	NewWriter(ctx context.Context, bucketName, objectName string) (io.WriteCloser, error)
}

type AmazonS3 interface {
	NewReader(ctx context.Context, region, bucket, key string) (io.ReadCloser, error)
	NewWriter(ctx context.Context, region, bucket, key string) (io.WriteCloser, error)
}
