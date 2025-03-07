---
page_title: "basis_vms Data Source - terraform-provider-bcc"
---
# basis_vms (Data Source)

Returns a list of Basis vms.

Get information about Vms in the Vdc for use in other resources.

Note: You can use the [`basis_vm`](Vm) data source to obtain metadata
about a single Vm if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_vms" "all_vms" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **vms** (List of Object) (see [below for nested schema](#nestedatt--vms))

<a id="nestedatt--vms"></a>
### Nested Schema for `vms`

Read-Only:

- **cpu** (Integer)
- **floating** (Boolean)
- **floating_ip** (String)
- **id** (String)
- **name** (String)
- **ram** (Integer)
- **template_id** (String)
- **template_name** (String)
- **power** (Boolean)
