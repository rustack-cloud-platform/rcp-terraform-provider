---
page_title: "basis_vdcs Data Source - terraform-provider-bcc"
---
# basis_vdcs (Data Source)

Get information about Vdcs in the Project for use in other resources.

Note: You can use the [`basis_vdc`](Vdc) data source to obtain metadata
about a single Vdc if you already know the `name` and unique `project_id`(optional) to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdcs" "all_vdcs" {
    project_id = data.basis_project.single_project.id
}

```

## Schema

### Required

- **project_id** (String) id of the Project

### Read-Only

- **vdcs** (List of Object) (see [below for nested schema](#nestedatt--vdcs))

<a id="nestedatt--vdcs"></a>
### Nested Schema for `vdcs`

Read-Only:

- **hypervisor** (String)
- **hypervisor_type** (String)
- **id** (String)
- **name** (String)
