package interfaces

import (
	"context"

	"github.com/secmon-as-code/locust/pkg/domain/model"
)

type UseCase interface {
	ValidateAzureEventGrid(ctx context.Context, callbackURL string) error

	HandleAzureCloudEvent(ctx context.Context, ev *model.CloudEventSchema) error
}
