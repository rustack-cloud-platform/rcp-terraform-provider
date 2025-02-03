---
page_title: "basis_lbaas Resource - terraform-provider-bcc"
---
# basis_lbaas (Resource)

Provides a Basis DNS record resource.

## Example Usage

```hcl
data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_network" "new_network" {
    vdc_id =  data.basis_vdc.single_vdc.id
    name = "New network"
}

resource "basis_lbaas" "lbaas" {
    vdc_id = data.basis_project.single_vdc.id
    name = "lbaas"
    port{
        network_id = data.basis_network.new_network.id
    }
    tags = ["created_by:terraform"]
}

```

## Schema

### Required

- **vdc_id** (String) id of Vdc
- **name** (String) name of LoadBalancer
- **Port** (String) parameter that specifies which network will be connected to LoadBalancer  (see [below for nested schema](#nestedblock--port))


### Optional

- **floating** (Boolean) enable floating ip for the LoadBalancer.
- **tags** (Toset, String) list of Tags added to the LoadBalancer.
- **timeouts** (Block, Optional)

<a id="nestedblock--port"></a>
### Nested Schema for `port`

Required:

- **network_id** (String) id of the Network

Optional:

- **ip_address** (String) ip address of port
