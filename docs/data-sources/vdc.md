---
page_title: "basis_vdc Data Source - terraform-provider-bcc"
---
# basis_vdc (Data Source)

Get information about a Vdc for use in other resources. 
This is useful if you need to utilize any of the Vdc's data and Vdc not managed by Terraform.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_vdc" "single_vdc2" {
    name = "Terraform VDC"
    # or
    id = "id"
}

```

## Schema

### Required

- **name** (String) name of the vdc `or` **id** (String) id of the vdc

### Optional

- **project_id** (String) id of the Project

### Read-Only

- **hypervisor** (String) name of the Hypervisor
- **hypervisor_type** (String) type of the Hypervisor
