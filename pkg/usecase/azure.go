package usecase

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/locust/pkg/domain/context/logging"
	"github.com/secmon-as-code/locust/pkg/domain/model"
)

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

func (x *UseCase) HandleAzureCloudEvent(ctx context.Context, ev *model.CloudEventSchema) error {
	logger := logging.From(ctx)
	logger.Debug("Handle Azure CloudEvent", "event", ev)

	// Example:
	// "/blobServices/default/containers/xxx-logs/blobs/tenantId=1yyyyy-yyyy-yyyy-yyyyyyyyyyyy/y=2024/m=08/d=25/h=23/m=00/PT1H.json"
	subject := strings.Split(ev.Subject, "/")
	if len(subject) < 6 {
		return goerr.New("Invalid Azure EventGrid message").With("subject", ev.Subject)
	}
	if subject[1] != "blobServices" || subject[3] != "containers" || subject[5] != "blobs" {
		return goerr.New("Invalid Azure EventGrid message").With("subject", ev.Subject)
	}

	// Example:
	// "/subscriptions/xxxx-xxxx-xxxx-xxxx/resourceGroups/xxxx/providers/Microsoft.Storage/storageAccounts/xxxx"
	source := strings.Split(ev.Source, "/")
	if len(source) != 9 {
		return goerr.New("invalid Azure EventGrid message").With("source", ev.Source)
	}
	if source[1] != "subscriptions" || source[3] != "resourceGroups" || source[5] != "providers" || source[7] != "storageAccounts" {
		return goerr.New("invalid Azure EventGrid message").With("source", ev.Source)
	}

	input := &model.RouteInput{
		AzureBlobStorage: &model.AzureBlobStorageEvent{
			Event: *ev,
			Object: model.AzureBlobStorageObject{
				StorageAccount: source[8],
				Container:      subject[4],
				BlobName:       strings.Join(subject[6:], "/"),
				ContentLength:  ev.Data.ContentLength,
				ContentType:    ev.Data.ContentType,
			},
		},
	}

	if err := x.Route(ctx, input); err != nil {
		return goerr.Wrap(err, "failed to emit route").With("input", input)
	}

	return nil
}
