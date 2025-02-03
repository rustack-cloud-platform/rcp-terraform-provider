---
page_title: "basis_disks Data Source - terraform-provider-bcc"
---
# basis_disks (Data Source)

Get information about list of Disks in the Vdc for use in other resources.

Note: You can use the [`basis_storage_profile`](Disk) data source to obtain metadata
about a single Disk if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_disks" "all_disks" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **disks** (List of Object) (see [below for nested schema](#nestedatt--disks))

<a id="nestedatt--disks"></a>
### Nested Schema for `disks`

Read-Only:

- **id** (String)
- **name** (String)
- **size** (Integer)
- **storage_profile_id** (String)
- **storage_profile_name** (String)
