---
page_title: "basis_dnss Data Source - terraform-provider-bcc"
---
# basis_dnss (Data Source)

Get information about Dnss in the Project for use in other resources.

Note: You can use the [`basis_dns`](Dns) data source to obtain metadata
about a single Dns if you already know the `name` and unique `project_id` to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_dnss" "dns" {
    project_id = data.basis_project.single_project.id
}


```

## Schema

### Required

- **project_id** (String) id of the Project

### Read-Only

- **dnss** (List of Object) (see [below for nested schema](#nestedatt--dnss))

<a id="nestedatt--dnss"></a>
### Nested Schema for `dnss`

Read-Only:

- **id** (String)
- **name** (String)
