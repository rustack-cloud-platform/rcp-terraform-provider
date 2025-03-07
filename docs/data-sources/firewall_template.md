---
page_title: "basis_firewall_template Data Source - terraform-provider-bcc"
---
# basis_firewall_template (Data Source)

Get information about a Firewall Template for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_firewall_template" "single_template" {
    vdc_id = data.basis_vdc.single_vdc.id
    
    name = "Разрешить Web"
    # or
    id = "id"
}

```
## Schema

### Required

- **name** (String) name of the Template `or` **id** (String) id of the Template
- **vdc_id** (String) id of the VDC

### Read-Only

- **id** (String) id of the Template
