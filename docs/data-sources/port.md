---
page_title: "basis_port Data Source - terraform-provider-bcc"
---
# basis_port (Data Source)

Get information about a Port for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_port" "port" {
    vdc_id = data.basis_vdc.single_vdc.id
    ip_address = "0.0.0.0"
    id = "00000000-0000-0000-0000-000000000000"
}

```
## Schema

### Required

- **vdc_id** (String) id of the VDC
- **ip_address** (String) ip_address of the Port
- **id** (String) id of the Port

If both fields are specified (ip_address , id) search will be carried out by **id**

### Read-Only

- **network_id** (String) id of the Network
- **firewall_templates** (List of String) list of firewall rule ids of the Port

