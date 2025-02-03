---
page_title: "basis_s3_storage_bucket Resource - terraform-provider-bcc"
---
# basis_s3_storage_bucket (Resource)

This data source provides creating and deleting s3_bucket. You should have a s3_storage to create a s3 storage bucket.

## Example Usage

```hcl 
data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_s3_storage" "s3_storage" {
    project_id = resource.basis_project.project1.id
    name = "s3_storage"
}

resource "basis_s3_storage_bucket" "bucket" {
    s3_storage_id=data.basis_s3_storage.s3_storage.id
    name ="Bucket-"
}
```

## Schema

### Required

- **name** (String) name of the Vm
- **s3_storage_id** (String) id of the S3 Storage

### Read-Only

- **id** (String) The ID of this resource.
- **external_name** (String) external_name for the s3 bucket.
