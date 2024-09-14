package gcs_test

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-mizutani/gt"
	"github.com/secmon-as-code/nydus/pkg/adapter/gcs"
)

func TestIntegration(t *testing.T) {
	bucketName, ok := os.LookupEnv("TEST_GOOGLE_CLOUD_STORAGE_BUCKET_NAME")
	if !ok {
		t.Skip("Skip integration test")
	}
	client, err := gcs.New()
	gt.NoError(t, err)

	ctx := context.Background()
	objectKey := time.Now().Format("nydus-test/2006/01/02/15/")
	objectName := objectKey + uuid.NewString() + ".txt"

	// Write object
	w, err := client.NewWriter(ctx, bucketName, objectName)
	gt.NoError(t, err)
	gt.R1(w.Write([]byte("timeless words"))).NoError(t)
	gt.NoError(t, w.Close())

	// Read object
	r, err := client.NewReader(ctx, bucketName, objectName)
	gt.NoError(t, err)
	buf := gt.R1(io.ReadAll(r)).NoError(t)
	gt.NoError(t, r.Close())
	gt.Equal(t, string(buf), "timeless words")
}
