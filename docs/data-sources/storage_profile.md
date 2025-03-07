---
page_title: "basis_storage_profile Data Source - terraform-provider-bcc"
---
# basis_storage_profile (Data Source)

Get information about a Storage Profile for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_storage_profile" "single_storage_profile" {
    vdc_id = data.basis_vdc.single_vdc.id

    name = "ssd"
    # or
    id = "id"
}

```
## Schema

### Required

- **name** (String) Name of the storage profile `or` **id** (String) id of the storage profile
- **vdc_id** (String) id of the VDC
