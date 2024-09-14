package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/m-mizutani/goerr"
)

type Client struct {
	cred aws.CredentialsProvider
}

type Option func(*Client)

func WithCredentials(cred aws.CredentialsProvider) Option {
	return func(c *Client) {
		c.cred = cred
	}
}

func New(options ...Option) (*Client, error) {
	c := &Client{}

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

func (x *Client) NewReader(ctx context.Context, region, bucket, key string) (io.ReadCloser, error) {
	s3Client := s3.NewFromConfig(aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(nil),
	})

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	output, err := s3Client.GetObject(ctx, input)
	if err != nil {
		return nil, goerr.Wrap(err, "fail to get object").With("bucket", bucket).With("key", key)
	}
	return output.Body, nil
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

func (x *Client) NewWriter(ctx context.Context, region, bucket, key string) (io.WriteCloser, error) {
	s3Client := s3.NewFromConfig(aws.Config{
		Region:      region,
		Credentials: x.cred,
	})

	errCh := make(chan error, 1)
	r, w := io.Pipe()

	writer := &pipeWriter{
		w:     w,
		errCh: errCh,
	}

	input := &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   r,
	}

	go func() {
		defer close(errCh)
		if _, err := s3Client.PutObject(ctx, input); err != nil {
			errCh <- goerr.Wrap(err, "fail to put object").With("bucket", bucket).With("key", key)
			return
		}

		if err := r.Close(); err != nil {
			errCh <- goerr.Wrap(err, "fail to close reader").With("bucket", bucket).With("key", key)
		}
	}()

	return writer, nil
}
