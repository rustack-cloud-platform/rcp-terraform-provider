---
page_title: "basis_kubernetes Resource - terraform-provider-bcc"
---
# basis_kubernetes (Resource)

This data source provides creating and deleting kubernetes. You should have a vdc to create a kubernetes.
#
- After creation fields: `node_ram`, `node_cpu`, `node_disk_size`, `node_storage_profile_id`, `user_public_key_id`. 
- Will be used in update if field `nodes_count` has changed. Changes apply only to the fresh node

## Example Usage

```hcl 
data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_network" "service_network" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "Сеть"
}

data "basis_storage_profile" "ssd" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "ssd"
}

data "basis_account" "me"{}

data "basis_kubernetes_template" "kubernetes_template"{
    name = "Kubernetes 1.22.1"
    vdc_id = data.basis_vdc.single_vdc.id
    
}

data "basis_pub_key" "key"{
    name = "test"
    account_id = data.basis_account.me.id
}

data "basis_platform" "pl"{
    vdc_id = data.basis_vdc.single_vdc.id
    name = "Intel Cascade Lake"
    
}

resource "basis_kubernetes" "k8s" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "test"
    node_ram = 3
    node_cpu = 3
    platform = data.basis_platform.pl.id
    template_id = data.basis_kubernetes_template.kubernetes_template.id
    nodes_count = 2
    node_disk_size = 10
    node_storage_profile_id = data.basis_storage_profile.ssd.id
    user_public_key_id = data.basis_pub_key.key.id
    floating = true
    tags = ["created_by:terraform"]
}

output "dashboard_url" {
        value = resource.basis_kubernetes.k8s.dashboard_url
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC
- **name** (String) name of the Kubernetes
- **node_cpu** (Integer) the number virtual cpus of the Vm
- **node_ram** (Integer) memory of the Vm in gigabytes
- **template_id** (String) id of the Template
- **platform** (String) id of the Template `(this field only for vmware hypervisor)`
- **nodes_count** (Integer) id of the Template
- **node_disk_size** (Integer) Size of disk in Kubernetes node
- **node_storage_profile_id** (String) Storage profile of disk in Kubernetes node
- **user_public_key_id** (String) key for communicating between Kubernetes nodes

### Optional

- **floating** (Boolean) enable floating ip for the Kubernetes
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **tags** (Toset, String) list of Tags added to the Kubernetes.


### Read-Only

- **floating_ip** (String) floating ip for the Vm. May be omitted
- **id** (String) The ID of this resource.
- **dashboard_url** (String) URL to access kubernetes dashboard

## Getting information about kubernetes

### Get dashboard url
- *This block will print dashboard_url in console*
```
    output "dashboard_url" {
        value = resource.basis_kubernetes.k8s.dashboard_url
    }
```
### Get kubectl config
- *When kubernetes is created, the kubectl configuration will appears in workdir wolder*
