---
page_title: "basis_port Resource - terraform-provider-bcc"
---
# basis_port (Resource)

Provides a Basis port which can be attached to a VM and Router in order to provide connection with network.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_firewall_template" "allow_default" {
    vdc_id = resource.basis_vdc.vdc1.id
    name = "Разрешить входящие"
}


resource "basis_network" "network" {
    vdc_id = resource.basis_vdc.single_vdc.id
    name = "network"

    subnets {
        cidr = "10.20.3.0/24"
        dhcp = true
        gateway = "10.20.3.1"
        start_ip = "10.20.3.2"
        end_ip = "10.20.3.254"
        dns = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
    }
}

resource "basis_port" "router_port" {
    vdc_id = resource.basis_vdc.single_vdc.id

    network_id = resource.basis_network.network.id
    ip_address = "199.199.199.199"
    firewall_templates = [data.basis_firewall_template.allow_default.id]
    tags = ["created_by:terraform"]
}
```

## Schema

### Required

- **network_id** String) id of the Network
- **vdc_id** (String) id of the VDC

### Optional

- **firewall_templates** (List of String) list of firewall rule ids of the Port
- **ip_address** (String) ip address of port
- **tags** (Toset, String) list of Tags added to the Port.

### Read-Only

- **id** (String) id of the Port
