# Nydus

Cross-Cloud Platform Tool for Event-Driven Object Data Transfer.

![overview](https://github.com/user-attachments/assets/514b04ce-7ca7-4f68-830f-b94ca54f1d87)

The `nydus` copies object data between cloud storage services in an event-driven manner. It can receive notifications from cloud storage services and transfer object data between them. When an object is created, updated, or some action in a source storage service, `nydus` will automatically transfer the object to the destination storage service.

The name "nydus" comes from the [Nydus Network](https://starcraft.fandom.com/wiki/Nydus_network) in StarCraft, which is a network of tunnels that allows units to travel between locations.

## Use Cases

- **Backup data from one cloud storage service to another**. For example, coping backup data of critical database for your business into another cloud storage service for disaster recovery.
- **Centralized data management**. For example, copying data from multiple cloud storage services into a single cloud storage service for centralized data management. Some services can dump data into a specific cloud storage service, and `nydus` can transfer the data to the centralized cloud storage service.

## How it works

`nydus` is a HTTP server that listens to the events from cloud storage services. When an event is received, `nydus` will transfer the object data from the source storage service to the destination storage service.

Overview of the data transfer process:

1. `nydus` listens to the events from the source storage service as HTTP server.
    - Amazon S3 can send events via SNS (Simple Notification Service).
    - Google Cloud Storage can send via Pub/Sub.
    - Azure Blob Storage can send via Event Grid.
2. When an event is received, `nydus` parse event data and evaluate it with [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) policy.
3. If the result has "route" that describes the destination storage service, `nydus` will transfer the object data to the destination storage service.

## Getting Started

### Prerequisites

- For Google Cloud: You need Service Account to access Google Cloud Storage.
- For Azure: You need App to access Azure Blob Storage.
- For AWS: You need IAM Service Account to access Amazon S3.

### Write a Rego policy

Write a Rego policy that describes the routing rules for the object data transfer. The policy should return the destination storage service and the destination bucket name.

Here is an example of the Rego policy that routes the object data from Google Cloud Storage to Azure Blob Storage:

```rego
package route

gcs[dst] {
    dst := {
        "bucket": "my-backup-bucket",
        "name": sprintf("from-azure/%s/%s/%s", [
            input.abs.object.storage_account,
            input.abs.object.container,
            input.abs.object.blob_name,
        ]),
    }
}
```

See [How to write a Rego policy](#how-to-write-a-rego-policy) for more details.

### Creating your container image

Create a container image that contains the Rego policy and the `nydus` binary. `nydus` provides a Docker image that contains the `nydus` binary from the GitHub Container Registry. You can use the `nydus` image as a base image and copy the Rego policy into the image.

```Dockerfile
FROM ghcr.io/secmon-as-code/nydus:latest

# It assumes that the Rego policy is in the "policy" directory.
COPY policy /policy

ENV NYDUS_POLICY_DIR=/policy
ENV NYDUS_ADDR=:8080

ENTRYPOINT ["/nydus" , "serve"]
```

The `nydus` binary is located at `/nydus` in the container image. The Rego policy should be copied to the `/policy` directory in the container image.

Environment variables for the `nydus` binary:

- `NYDUS_POLICY_DIR` (required): The directory that contains the Rego policy files.
- `NYDUS_ADDR` (optional): The address that `nydus` listens to. The default value is `127.0.0.1:8080`. You need to set this environment variable to exposed binding address, such as `:8080` to listen to all interfaces.
- `NYDUS_LOG_LEVEL` (optional): The log level of `nydus`. The default value is `info`.
- `NYDUS_LOG_FORMAT` (optional): The log format of `nydus`. You can choose `console` or `json`. The default value is `json`.
- `NYDUS_ENABLE_GCS` (optional): Enable Google Cloud Storage client. It's required for both of download and upload an object. The default value is `false`. Following environment variables are required when `NYDUS_ENABLE_GCS` is `true`.
  - `NYDUS_GCS_CREDENTIAL_FILE` (optional): The path to the Google Cloud Service Account credential file. It's basically not needed when the application is running on Google Cloud Platform.
- `NYDUS_ENABLE_AZURE` (optional): Enable Azure Blob Storage client. It's required for both of download and upload an object. The default value is `false`. Following environment variables are required when `NYDUS_ENABLE_AZURE` is `true`.
  - `NYDUS_AZURE_TENANT_ID` (required): The Azure Tenant ID.
  - `NYDUS_AZURE_CLIENT_ID` (required): The Azure Client ID for the App.
  - `NYDUS_AZURE_CLIENT_SECRET` (required): The Azure Client Secret for the App.

### Deploy your container image

Deploy the container image to your container platform, such as Kubernetes, Docker, or any other container platform. We recommend using [Cloud Run](https://cloud.google.com/run?hl=en) on Google Cloud Platform, as it is a serverless container platform that can scale automatically.

## How to write a Rego policy

Please refer to the [Open Policy Agent documentation](https://www.openpolicyagent.org/docs/latest/policy-language/) for more details about the Rego policy language.

### Rego package name

The Rego policy should return the destination storage service information, such as the destination bucket name and the object path in the destination bucket. The policy should be written in the `route` package. That means the policy file should start with `package route`.

### Input data

The input data for the Rego policy is the event data from the source storage service. The event data is parsed by `nydus` and passed to the Rego policy as the `input` variable.

The `input` variable has the following structure:

- `abs`: The abstracted event data that is common to all cloud storage services.
    - `object`: The object data.
        - `storage_account`: The storage account name.
        - `container`: The container name.
        - `blob_name`: The blob name.
        - `size`: The object size.
        - `content_type`: The object content type.
        - `etag`: The object ETag.
    - `event`: This field contains original Azure Event Grid notification data. See [Azure Event Grid schema](https://docs.microsoft.com/en-us/azure/event-grid/event-schema-blob-storage?tabs=event-grid) for more details.
- `gcs`: To be supported soon.
- `s3`: To be supported soon.

### Output data

The Rego policy should return the destination storage service information as a set. The set should contain the destination bucket name and the object path in the destination bucket.

- `gcs`: The destination storage service is Google Cloud Storage. The variable must be [Set](https://www.openpolicyagent.org/docs/latest/policy-language/#sets) type and contain the following fields:
    - `bucket`: The destination bucket name.
    - `name`: The object path in the destination bucket.
- `s3`: To be supported soon.
- `abs`: To be supported soon.

## License

Apache License 2.0
