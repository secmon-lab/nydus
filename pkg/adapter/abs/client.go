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
