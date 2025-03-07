---
page_title: "basis_storage_profiles Data Source - terraform-provider-bcc"
---
# basis_storage_profiles (Data Source)

Get information about Storage Profiles in the Vdc for use in other resources.

Note: You can use the [`basis_storage_profile`](Storage Profile) data source to obtain metadata
about a single Storage Profile if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_storage_profiles" "all_storage_profiles" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **storage_profiles** (List of Object) (see [below for nested schema](#nestedatt--storage_profiles))

<a id="nestedatt--storage_profiles"></a>
### Nested Schema for `storage_profiles`

Read-Only:

- **id** (String)
- **name** (String)
