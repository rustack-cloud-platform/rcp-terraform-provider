---
page_title: "basis_lbaas Data Source - terraform-provider-bcc"
---
# basis_lbaas (Data Source)

Get information about Basis lbaas.

Get information about lbaas in the Vdc for use in other resources.

**Note:** This data source returns a single Lbaas. When specifying a `name`, an
error is triggered if more than one lbaas is found.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_lbaas" "test" {
    vdc_id = data.basis_vdc.single_vdc.id
    
    name = "test"
    # or
    id = "id"
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC
- **name** (String) name of the LoadBalancer `or` **id** (String) id of the LoadBalancer

### Read-Only

- **lbaas** (List of Object) (see [below for nested schema](#nestedatt--lbaas))

<a id="nestedatt--lbaas"></a>
### Nested Schema for `lbaas`

Read-Only:

- **floating** (Boolean)
- **name** (String)
