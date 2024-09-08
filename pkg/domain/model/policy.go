package model

type RouteInput struct {
	AzureBlobStorage *AzureBlobStorageEvent `json:"abs"`
}

type RouteOutput struct {
	AzureBlobStorage   []AzureBlobStorageObject   `json:"abs"`
	GoogleCloudStorage []GoogleCloudStorageObject `json:"gcs"`
	AmazonS3Storage    []AmazonS3Object           `json:"s3"`
}
