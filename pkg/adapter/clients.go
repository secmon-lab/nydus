package adapter

import (
	"net/http"

	"github.com/m-mizutani/opac"
	"github.com/secmon-as-code/nydus/pkg/domain/interfaces"
)

type Clients struct {
	httpClient interfaces.HTTPClient
	query      *opac.Client

	gcsClient interfaces.GoogleCloudStorage
	absClient interfaces.AzureBlobStorage
	s3Client  interfaces.AmazonS3
}

func (x *Clients) HTTPClient() interfaces.HTTPClient { return x.httpClient }
func (x *Clients) Query() *opac.Client               { return x.query }
func (x *Clients) GoogleCloudStorage() interfaces.GoogleCloudStorage {
	return x.gcsClient
}
func (x *Clients) AzureBlobStorage() interfaces.AzureBlobStorage {
	return x.absClient
}
func (x *Clients) AmazonS3() interfaces.AmazonS3 { return x.s3Client }

func New(options ...Option) *Clients {
	clients := &Clients{
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(clients)
	}

	return clients
}

type Option func(*Clients)

func WithHTTPClient(httpClient interfaces.HTTPClient) Option {
	return func(c *Clients) {
		c.httpClient = httpClient
	}
}

func WithPolicy(query *opac.Client) Option {
	return func(c *Clients) {
		c.query = query
	}
}

func WithGoogleCloudStorage(client interfaces.GoogleCloudStorage) Option {
	return func(c *Clients) {
		c.gcsClient = client
	}
}

func WithAzureBlobStorage(client interfaces.AzureBlobStorage) Option {
	return func(c *Clients) {
		c.absClient = client
	}
}

func WithAmazonS3(client interfaces.AmazonS3) Option {
	return func(c *Clients) {
		c.s3Client = client
	}
}
