package adapter

import (
	"net/http"

	"github.com/m-mizutani/opac"
	"github.com/secmon-as-code/locust/pkg/domain/interfaces"
)

type Clients struct {
	httpClient interfaces.HTTPClient
	query      opac.Client
}

func (x *Clients) HTTPClient() interfaces.HTTPClient { return x.httpClient }
func (x *Clients) Query() opac.Client                { return x.query }

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

func WithPolicy(query opac.Client) Option {
	return func(c *Clients) {
		c.query = query
	}
}
