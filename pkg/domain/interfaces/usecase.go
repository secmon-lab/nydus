package interfaces

import (
	"context"

	"github.com/secmon-as-code/nydus/pkg/domain/model"
)

type UseCase interface {
	ValidateAzureEventGrid(ctx context.Context, callbackURL string) error

	HandleAzureCloudEvent(ctx context.Context, ev *model.CloudEventSchema) error
}
