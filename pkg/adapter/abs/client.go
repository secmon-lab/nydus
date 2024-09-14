package abs

import (
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/m-mizutani/goerr"
)

type Client struct {
	cred *azidentity.ClientSecretCredential
}

func New(clientID, clientSecret, tenantID string) (*Client, error) {
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create azure client").With("clientID", clientID).With("tenantID", tenantID).With("clientSecret.length", len(clientSecret))
	}

	return &Client{
		cred: cred,
	}, nil
}

func (x *Client) NewReader(ctx context.Context, storageAccountName, containerName, blobName string) (io.ReadCloser, error) {
	accountUrl := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)

	serviceClient, err := azblob.NewClient(accountUrl, x.cred, nil)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create service client").With("accountUrl", accountUrl)
	}

	stream, err := serviceClient.DownloadStream(ctx, containerName, blobName, &azblob.DownloadStreamOptions{})
	if err != nil {
		return nil, goerr.Wrap(err, "fail to download stream").With("containerName", containerName).With("blobName", blobName).With("accountUrl", accountUrl)
	}

	return stream.Body, nil
}

type pipeWriter struct {
	w     io.WriteCloser
	errCh chan error
}

func (x *pipeWriter) Write(p []byte) (n int, err error) {
	return x.w.Write(p)
}

func (x *pipeWriter) Close() error {
	if err := x.w.Close(); err != nil {
		return err
	}

	return <-x.errCh
}

func (x *Client) NewWriter(ctx context.Context, storageAccountName, containerName, blobName string) (io.WriteCloser, error) {
	accountUrl := fmt.Sprintf("https://%s.blob.core.windows.net/", storageAccountName)

	serviceClient, err := azblob.NewClient(accountUrl, x.cred, nil)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create service client").With("accountUrl", accountUrl)
	}

	errCh := make(chan error)
	r, w := io.Pipe()

	writer := &pipeWriter{
		w:     w,
		errCh: errCh,
	}

	go func() {
		defer close(errCh)

		if _, err := serviceClient.UploadStream(ctx, containerName, blobName, r, nil); err != nil {
			errCh <- goerr.Wrap(err, "fail to create writer").With("containerName", containerName).With("blobName", blobName).With("accountUrl", accountUrl)
			return
		}

		if err := r.Close(); err != nil {
			errCh <- goerr.Wrap(err, "fail to close reader").With("containerName", containerName).With("blobName", blobName).With("accountUrl", accountUrl)
		}
	}()

	return writer, nil
}
