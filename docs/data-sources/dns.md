---
page_title: "basis_dns Data Source - terraform-provider-bcc"
---
# basis_dns (Data Source)

Get information about a Dns for use in other resources. 
This is useful if you need to utilize any of the Dns's data and dns not managed by Terraform.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_dns" "dns" {
    name="dns.teraform."
    # or
    id = "id"
    
    project_id = data.basis_project.single_project.id
}

```

## Schema

### Required

- **project_id** (String) id of the Project
- **name** (String) name of the dns zone `or` **id** (String) id of the dns zone


