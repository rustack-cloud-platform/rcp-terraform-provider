---
page_title: "basis_hypervisors Data Source - terraform-provider-bcc"
---
# basis_hypervisors (Data Source)

Get information about list of Hypervisors in the Project for use in other resources.

Note: You can use the [`basis_hypervisor`](Hypervisor) data source to obtain metadata
about a single Hypervisor if you already know the `name` and `project_id` to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_hypervisors" "all_hypervisors" {
    project_id = data.basis_project.single_project.id
}

```

## Schema

### Required

- **project_id** (String) id of the Project

### Read-Only

- **hypervisors** (List of Object) (see [below for nested schema](#nestedatt--hypervisors))

<a id="nestedatt--hypervisors"></a>
### Nested Schema for `hypervisors`

Read-Only:

- **id** (String)
- **name** (String)
- **type** (String)
