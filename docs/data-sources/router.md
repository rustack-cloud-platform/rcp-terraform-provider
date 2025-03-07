---
page_title: "basis_router Data Source - terraform-provider-bcc"
---
# basis_router (Data Source)

Get information about a Router for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id"
    name = "Terraform VDC"
}

data "basis_router" "single_Router" {
    vdc_id = data.basis_vdc.single_vdc.id
    
    name = "Terraform Router"
    # or
    id = "id"
}

```
## Schema

### Required

- **vdc_id** (String) id of the VDC
- **name** (String) name of the Router `or` **id** (String) id of the Router

