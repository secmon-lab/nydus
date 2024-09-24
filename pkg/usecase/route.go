package usecase

import (
	"context"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-as-code/nydus/pkg/adapter"
	"github.com/secmon-as-code/nydus/pkg/domain/context/logging"
	"github.com/secmon-as-code/nydus/pkg/domain/model"
)

func (x *UseCase) Route(ctx context.Context, input *model.RouteInput) error {
	var output model.RouteOutput

	logger := logging.From(ctx)
	logger.Debug("Route query", "input", input, "clients", x.clients)
	logger.Debug("[tmp] client", "clients", x.clients)

	if err := x.clients.Query().Query(ctx, "data.route", input, &output); err != nil {
		return goerr.Wrap(err, "failed to route query").With("input", input)
	}
	logger.Info("Route query result", "input", input, "output", output)

	for _, dst := range output.GoogleCloudStorage {
		logger.Debug("Route to Google Cloud Storage", "destination", dst)

		r, err := newReaderFromRouteInput(ctx, x.clients, input)
		if err != nil {
			return goerr.Wrap(err, "failed to create reader from route input").With("input", input)
		}
		defer r.Close()

		w, err := x.clients.GoogleCloudStorage().NewWriter(ctx, dst.Bucket, dst.Name)
		if err != nil {
			return goerr.Wrap(err, "failed to create writer to Google Cloud Storage").With("destination", dst)
		}
		defer w.Close()

		if _, err := io.Copy(w, r); err != nil {
			return goerr.Wrap(err, "failed to copy from reader to writer")
		}
	}

	return nil
}

func newReaderFromRouteInput(ctx context.Context, clients *adapter.Clients, input *model.RouteInput) (io.ReadCloser, error) {
	logging.From(ctx).Info("Create reader from route input", "input", input, "clients", clients)

	switch {
	case input.AzureBlobStorage != nil:
		return clients.AzureBlobStorage().NewReader(ctx,
			input.AzureBlobStorage.Object.StorageAccount,
			input.AzureBlobStorage.Object.Container,
			input.AzureBlobStorage.Object.BlobName,
		)
	case input.GoogleCloudStorage != nil:
		return clients.GoogleCloudStorage().NewReader(ctx,
			input.GoogleCloudStorage.Object.Bucket,
			input.GoogleCloudStorage.Object.Name,
		)

	case input.AmazonS3 != nil:
		return clients.AmazonS3().NewReader(ctx,
			input.AmazonS3.Object.Region,
			input.AmazonS3.Object.Bucket,
			input.AmazonS3.Object.Key,
		)
	default:
		return nil, goerr.New("unsupported route input")
	}
}
