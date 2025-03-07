---
page_title: "basis_kubernetess Data Source - terraform-provider-bcc"
---
# basis_kubernetess (Data Source)

Returns a list of Basis kubernetess.

Get information about all kubernetes clasters in the Vdc for use in other resources.

Note: You can use the [`basis_vm`](Kubernetess) data source to obtain metadata
about a single Kubernetess if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_kubernetess" "all_k8s" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **kubernetess** (List of Object) (see [below for nested schema](#nestedatt--kubernetess))

<a id="nestedatt--kubernetess"></a>
### Nested Schema for `kubernetess`

Read-Only:

- **id** (String)
- **name** (String)
- **node_cpu** (Integer)
- **node_ram** (Integer)
- **template_id** (String)
- **floating** (Boolean)
- **floating_ip** (String)
- **node_disk_size** (String)
- **nodes_count** (String)
- **user_public_key_id** (String)
- **node_storage_profile_id** (String)
- **vms** (String)
- **dashboard_url** (String)

## Getting information about kubernetes

### Get dashboard url

- *This block will print dashboard_url in console*
```
    output "dashboard_url" {
        value = data.basis_kubernetes.all_k8s[0].dashboard_url
    }
```

### Get kubectl config
- *When kubernetes is received, the kubectl configuration will appear in the workdir folder.*
