package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/m-mizutani/goerr"
	"google.golang.org/api/option"
)

// Client is a client for Google Cloud Storage
type Client struct {
	client  *storage.Client
	options []option.ClientOption
}

type Option func(*Client)

func New(options ...Option) (*Client, error) {
	c := &Client{}

	for _, opt := range options {
		opt(c)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, c.options...)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create GCS client")
	}

	c.client = client

	return c, nil
}

func WithGoogleAPIOption(opts ...option.ClientOption) Option {
	return func(c *Client) {
		c.options = opts
	}
}

func (x *Client) NewReader(ctx context.Context, bucket, object string) (io.ReadCloser, error) {
	reader, err := x.client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to create reader").With("bucket", bucket).With("object", object)
	}

	return reader, nil
}

func (x *Client) NewWriter(ctx context.Context, bucket, object string) (io.WriteCloser, error) {
	writer := x.client.Bucket(bucket).Object(object).NewWriter(ctx)
	return writer, nil
}
