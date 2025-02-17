package model

type RouteInput struct {
	AzureBlobStorage   *AzureBlobStorageEvent   `json:"abs"`
	GoogleCloudStorage *GoogleCloudStorageEvent `json:"gcs"`
	AmazonS3           *AmazonS3Event           `json:"s3"`
	Env                map[string]string        `json:"env"`
}

type RouteOutput struct {
	AzureBlobStorage   []AzureBlobStorageObject   `json:"abs"`
	GoogleCloudStorage []GoogleCloudStorageObject `json:"gcs"`
	AmazonS3Storage    []AmazonS3Object           `json:"s3"`
}
