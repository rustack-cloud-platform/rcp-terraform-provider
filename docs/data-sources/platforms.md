---
page_title: "basis_platforms Data Source - terraform-provider-bcc"
---
# basis_platforms (Data Source)
### `Only for Vmware Hypervisor`

Get information about Platforms in the Vdc for use in other resources.

Note: You can use the [`basis_platforms`](Platforms) data source to obtain metadata
about a single Platforms if you already know the `name` to retrieve.

## Example Usage

```hcl
data "basis_project" "single_project" {
    name = "Terraform Project"
}
data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}
data "basis_platforms" "platforms"{
    vdc_id = resource.basis_vdc.single_vdc.id
}
```

## Schema

### Read-Only

- **platforms** (List of Object) (see [below for nested schema](#nestedatt--projects))

<a id="nestedatt--platforms"></a>
### Nested Schema for `platforms`

Read-Only:

- **id** (String)
- **name** (String)