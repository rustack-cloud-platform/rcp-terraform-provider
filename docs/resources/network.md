---
page_title: "basis_network Resource - terraform-provider-bcc"
---
# basis_network (Resource)

Provides a Basis network to provide connections of two or more computers that are linked in order to share resources.

## Example Usage

```hcl
data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

resource "basis_network" "network2" {
    vdc_id = data.basis_vdc.single_vdc.id

    name = "Network 1"

    subnets {
        cidr = "10.20.1.0/24"
        dhcp = true
        gateway = "10.20.1.1"
        start_ip = "10.20.1.2"
        end_ip = "10.20.1.254"
        dns = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
    }
    tags = ["created_by:terraform"]
}
```

## Schema

### Required

- **name** (String) name of the Network
- **subnets** (Block List, Min: 1, Max: 1) (see [below for nested schema](#nestedblock--subnets))
- **vdc_id** (String) id of the VDC

### Optional

- **id** (String) The ID of this resource.
- **tags** (Toset, String) list of Tags added to the Network.
- **mtu** (Integer) maximum transmission unit for the Network

<a id="nestedblock--subnets"></a>
### Nested Schema for `subnets`

Required:

- **cidr** (String) cidr of the Subnet
- **dhcp** (Boolean) enable dhcp service of the Subnet
- **dns** (List of String) dns servers list
- **end_ip** (String) pool end ip of the Subnet
- **gateway** (String) gateway of the Subnet
- **start_ip** (String) pool start ip of the Subnet

Read-Only:

- **id** (String) id of the Subnet
