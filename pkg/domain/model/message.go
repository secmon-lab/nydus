package model

type StorageType string

const (
	AzureBlobStorage   StorageType = "abs"
	GoogleCloudStorage StorageType = "gcs"
	S3Storage          StorageType = "s3"
)

// AzureBlobStorageEvent is a struct for Azure Blob Storage event
type AzureBlobStorageEvent struct {
	Event  CloudEventSchema       `json:"event"`
	Object AzureBlobStorageObject `json:"object"`
}

type AzureBlobStorageObject struct {
	StorageAccount string `json:"storage_account"`
	Container      string `json:"container"`
	BlobName       string `json:"blob_name"`
	ContentLength  int64  `json:"contentLength"`
	ContentType    string `json:"contentType"`
}

// CloudEventSchema is a struct for Azure Event Grid CloudEvent schema
type CloudEventSchema struct {
	Data struct {
		API                string `json:"api"`
		BlobType           string `json:"blobType"`
		ClientRequestID    string `json:"clientRequestId"`
		ContentLength      int64  `json:"contentLength"`
		ContentType        string `json:"contentType"`
		ETag               string `json:"eTag"`
		RequestID          string `json:"requestId"`
		Sequencer          string `json:"sequencer"`
		StorageDiagnostics struct {
			BatchID string `json:"batchId"`
		} `json:"storageDiagnostics"`
		URL string `json:"url"`
	} `json:"data"`
	ID          string `json:"id"`
	Source      string `json:"source"`
	SpecVersion string `json:"specversion"`
	Subject     string `json:"subject"`
	Time        string `json:"time"`
	Type        string `json:"type"`
}

type GoogleCloudStorageEvent struct {
	Event  GooglePubSubEvent        `json:"event"`
	Object GoogleCloudStorageObject `json:"object"`
}

// GoogleCloudStorageObject is a struct for Google Cloud Storage object
type GoogleCloudStorageObject struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

type GooglePubSubEvent struct {
}

type AmazonS3Event struct {
	Event  AmazonSNSEvent `json:"event"`
	Object AmazonS3Object `json:"object"`
}

type AmazonSNSEvent struct {
}

type AmazonS3Object struct {
	Region string `json:"region"`
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}
