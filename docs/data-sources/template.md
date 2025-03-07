---
page_title: "basis_template Data Source - terraform-provider-bcc"
---
# basis_template (Data Source)

Get information about a Template for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_template" "single_template" {
    vdc_id = data.basis_vdc.single_vdc.id
    
    name = "Debian 10"
    # or
    id = "id"
}

```

## Schema

### Required

- **name** (String) name of the Template `or` **id** (String) id of the Template
- **vdc_id** (String) id of the VDC

### Read-Only

- **min_cpu** (Integer) minimum cpu required by the Template
- **min_disk** (Integer) minimum disk size in GB required by the Template
- **min_ram** (Float) minimum ram in GB required by the Template
