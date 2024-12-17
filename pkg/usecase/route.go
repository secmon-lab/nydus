package usecase

import (
	"context"
	"io"

	"github.com/m-mizutani/goerr"
	"github.com/secmon-lab/nydus/pkg/adapter"
	"github.com/secmon-lab/nydus/pkg/domain/context/logging"
	"github.com/secmon-lab/nydus/pkg/domain/model"
)

func (x *UseCase) Route(ctx context.Context, input *model.RouteInput) (err error) {
	var output model.RouteOutput

	logger := logging.From(ctx)
	logger.Debug("Route query", "input", input)
	if err := x.clients.Query().Query(ctx, "data.route", input, &output); err != nil {
		return goerr.Wrap(err, "failed to route query").With("input", input)
	}
	logger.Info("Route query result", "input", input, "output", output)

	for _, dst := range output.GoogleCloudStorage {
		logger.Debug("Route to Google Cloud Storage", "destination", dst)
		gcs := x.clients.GoogleCloudStorage()
		if gcs == nil {
			return goerr.New("Google Cloud Storage is not enabled").With("destination", dst).With("input", input)
		}

		r, err := newReaderFromRouteInput(ctx, x.clients, input)
		if err != nil {
			return goerr.Wrap(err, "failed to create reader from route input").With("input", input)
		}
		defer r.Close()

		w, err := x.clients.GoogleCloudStorage().NewWriter(ctx, dst.Bucket, dst.Name)
		if err != nil {
			return goerr.Wrap(err, "failed to create writer to Google Cloud Storage").With("destination", dst)
		}
		defer func() {
			if err = w.Close(); err != nil {
				logger.Warn("Failed to close writer", "destination", dst, "error", err)
			}
		}()

		n, err := io.Copy(w, r)
		if err != nil {
			return goerr.Wrap(err, "failed to copy from reader to writer")
		}

		logger.Info("Copied from reader to writer", "destination", dst, "bytes", n)
	}

	return nil
}

func newReaderFromRouteInput(ctx context.Context, clients *adapter.Clients, input *model.RouteInput) (io.ReadCloser, error) {
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
