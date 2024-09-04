package usecase

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/locust/pkg/adapter"
	"github.com/secmon-as-code/locust/pkg/domain/context/logging"
)

type UseCase struct {
	clients *adapter.Clients
}

func New(clients *adapter.Clients) *UseCase {
	return &UseCase{
		clients: clients,
	}
}

func (x *UseCase) ValidateAzureEventGrid(ctx context.Context, callbackURL string) error {
	reqURL, err := url.Parse(callbackURL)
	if err != nil {
		return goerr.Wrap(err, "Failed to parse callbackURL")
	}

	// Example:
	// https://rp-japaneast.eventgrid.azure.net:553/eventsubscriptions/my-topic/validate?id=XXXXXXXXX&t=2024-08-25T23:16:11.8746191Z&apiVersion=2023-12-15-preview&token=XXXXXXX%3d
	if reqURL.Scheme != "https" ||
		!strings.HasSuffix(reqURL.Hostname(), ".eventgrid.azure.net") {
		return goerr.New("Webhook-Request-Callback is invalid")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, callbackURL, nil)
	if err != nil {
		return goerr.Wrap(err, "failed to create HTTP request").With("callbackURL", callbackURL)
	}

	resp, err := x.clients.HTTPClient().Do(req)
	if err != nil {
		return goerr.Wrap(err, "failed to send HTTP request").With("callbackURL", callbackURL)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return goerr.New("callbackURL response is not OK").With("statusCode", resp.StatusCode).With("body", string(body)).With("callbackURL", callbackURL)
	}

	logging.From(ctx).Info("Successfully validated Azure EventGrid", "callbackURL", callbackURL)

	return nil
}
