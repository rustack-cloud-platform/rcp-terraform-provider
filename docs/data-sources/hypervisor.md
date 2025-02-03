---
page_title: "basis_hypervisor Data Source - terraform-provider-bcc"
---
# basis_hypervisor (Data Source)

Get information about a Hypervisor for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_hypervisor" "single_hypervisor" {
    project_id = data.basis_project.single_project.id
    
    name = "VMWARE"
    # or
    id ="id"
}

```

## Schema

### Required

- **name** (String) name of the Hypervisor `or` **id** (String) id of the Hypervisor
- **project_id** (String) id of the Project

### Read-Only

- **type** (String) type of the Hypervisor
