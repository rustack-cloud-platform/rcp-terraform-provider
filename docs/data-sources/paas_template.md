---
page_title: "basis_paas_template Data Source - terraform-provider-bcc"
---
# basis_paas_template (Data Source)

Get information about a PaaS Service Template for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_paas_template" "paas_template" {
  id = 1
  project_id = data.basis_project.single_project.id
}
```
## Schema

### Required

- **project_id** (String) id of Project
- **id** (String) id of PaaS Service Template

### Read-Only

- **name** (String) name of PaaS Service Template
