---
page_title: "basis_platform Data Source - terraform-provider-bcc"
---
# basis_platform (Data Source)
### `Only for Vmware Hypervisor`
Get information about a Platform for use in other resources. 

## Example Usage

```hcl
data "basis_project" "single_project" {
    name = "Terraform Project"
}
data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}
data "basis_platform" "platform"{
    vdc_id = resource.basis_vdc.single_vdc.id
    name = "Intel Cascade Lake"
    # or
    id = ""
}
```
## Schema

### Required

- **name** (String) name of the Platform `or` **id** (String) id of the Platform
