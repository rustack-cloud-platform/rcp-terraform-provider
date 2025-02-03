---
page_title: "basis_lbaas_pool Resource - terraform-provider-bcc"
---
# basis_lbaas_pool (Resource)

Provides a Basis Lbaas pool resource.

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

data "basis_template" "debian10" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "Debian 10"
}

data "basis_firewall_template" "allow_default" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "По-умолчанию"
}

data "basis_lbaas" "lbaas" {
    vdc_id = data.basis_project.single_vdc.id
    name = "lbaas"
}

data "basis_port" "vm_port" {
    vdc_id = resource.basis_vdc.single_vdc.id

    network_id = resource.basis_network.new_network.id
    firewall_templates = [data.basis_firewall_template.allow_default.id]
}

data "basis_vm" "vm" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "Server 1"
}

resource "basis_lbaas_pool" "pool" {
    lbaas_id = data.basis_lbaas.lbaas.id
    connlimit = 65536
    method = "ROUND_ROBIN"
    port = 2
    protocol = "TCP"
    member {
        port = 2
        weight = 1
        vm_id = data.basis_vm.vm.id
    }
    
    depends_on = [basis_vm.vm]
}

```

## Schema

### Required

- **lbaas_id** (String) id of LoadBalancer
- **port** (Integer) port of LoadBalancerPool
- **member** (String) parameter that specifies which network will be connected to LoadBalancer  (see [below for nested schema](#nestedblock--member))


### Optional

- **method** (String) method of LoadBalancerPool 
> Can be chosen ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP
- **protocol** (String) method of LoadBalancerPool
> Can be chosen TCP, HTTP, HTTPS
- **connlimit** (Integer) connlimit of LoadBalancerPool
- **timeouts** (Block, Optional)

<a id="nestedblock--member"></a>
### Nested Schema for `member`

Required:

- **port** (Integer) id of the Network
- **vm_id** (String) id of the Network

Optional:

- **weight** (Integer) id of the Network
