---
page_title: "basis_network Data Source - terraform-provider-bcc"
---
# basis_network (Data Source)

Get information about a Network for use in other resources. 

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_network" "single_network" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "Сеть 1"
    # or
    id = "id"
}

```
## Schema

### Required

- **name** (String) name of the Network `or` **id** (String) id of the Network
- **vdc_id** (String) id of the VDC

### Read-Only


- **mtu** (Integer) maximum transmission unit for the Network
- **subnets** (List of Object) list of subnets (see [below for nested schema](#nestedatt--subnets))

<a id="nestedatt--subnets"></a>
### Nested Schema for `subnets`

Read-Only:

- **cidr** (String)
- **dhcp** (Boolean)
- **dns** (List of String)
- **end_ip** (String)
- **gateway** (String)
- **id** (String)
- **start_ip** (String)
