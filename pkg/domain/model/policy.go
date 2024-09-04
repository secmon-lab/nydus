package model

type RouteInput struct {
	BlobStorage *AzureBlobStorageEvent `json:"blob_storage"`
}

type RouteOutput struct {
	AzureBlobStorage   []AzureBlobStorageObject   `json:"abs"`
	GoogleCloudStorage []GoogleCloudStorageObject `json:"gcs"`
	AmazonS3Storage    []AmazonS3Object           `json:"s3"`
}
