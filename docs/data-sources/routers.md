---
page_title: "basis_router Data Source - terraform-provider-bcc"
---
# basis_router (Data Source)

Get information about a Routers for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id"
    name = "Terraform VDC"
}

data "basis_routers" "vdc_routers" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```
## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **routers** (List of Object) (see [below for nested schema](#nestedatt--router))

<a id="nestedatt--router"></a>
### Nested Schema for `router`

Read-Only:

- **id** (String)
- **name** (String)
