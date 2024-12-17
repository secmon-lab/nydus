package usecase_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/nydus/pkg/adapter"
	"github.com/secmon-lab/nydus/pkg/usecase"
)

type mockHTTPClient struct {
	requests []*http.Request
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.requests = append(c.requests, req)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("Webhook successfully validated as a subscription endpoint."))),
	}, nil
}

func TestAzureValidation(t *testing.T) {
	const testURL = "https://rp-japaneast.eventgrid.azure.net:553/eventsubscriptions/xxxxxx/validate?id=XXXXX-XXXXXX-XXXXX&t=2024-12-17T08:20:31.2630520Z&apiVersion=2023-12-15-preview&token=Z%2f%oiujgoafasiodjfaposdijfasd%3d"

	mock := &mockHTTPClient{}
	uc := usecase.New(adapter.New(adapter.WithHTTPClient(mock)))

	ctx := context.Background()
	err := uc.ValidateAzureEventGrid(ctx, testURL)
	gt.NoError(t, err)

	gt.A(t, mock.requests).Length(1).At(0, func(t testing.TB, v *http.Request) {
		gt.Equal(t, v.URL.String(), testURL)
	})
}
