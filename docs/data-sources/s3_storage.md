---
page_title: "basis_s3_storage Data Source - terraform-provider-bcc"
---
# basis_s3_storage (Data Source)

Get information about a S3Storage for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_s3_storage" "s3_storage" {
    project_id = resource.basis_project.single_project.id
    
    name = "s3_storage"
    # or
    id = "id"
}

```

## Schema

### Required

- **project_id** (String) id of the project
- **name** (String) name of the S3Storage `or` **id** (String) id of the S3Storage

### Read-Only

- **backend** (String) backend for access to s3 (`minio` or `netapp`)
- **client_endpoint** (String) url for connecting to s3"
- **access_key** (String) access_key for access to s3
- **secret_key** (String) secret_key for access to s3