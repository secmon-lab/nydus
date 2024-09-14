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

### Output data

The Rego policy should return the destination storage service information as a set. The set should contain the destination bucket name and the object path in the destination bucket.

- `gcs`: The destination storage service is Google Cloud Storage. The variable must be [Set](https://www.openpolicyagent.org/docs/latest/policy-language/#sets) type and contain the following fields:
    - `bucket`: The destination bucket name.
    - `name`: The object path in the destination bucket.

## License

Apache License 2.0
