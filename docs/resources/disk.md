---
page_title: "basis_disk Resource - terraform-provider-bcc"
---
# basis_disk (Resource)

Provides a Basis disk volume which can be attached to a VM in order to provide expanded storage.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_storage_profile" "single_storage_profile" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "sas"
}

resource "basis_disk" "disk2" {
    vdc_id = data.basis_vdc.single_vdc.id

    name = "Disk 1"
    storage_profile_id = data.basis_storage_profile.single_storage_profile.id
    size = 1
    tags = ["created_by:terraform"]
}
```

## Schema

### Required

- **name** (String) name of the Disk
- **size** (Integer) the size of the Disk in gigabytes
- **storage_profile_id** (String) Id of the storage profile
- **vdc_id** (String) id of the VDC

### Optional

- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- **tags** (Toset, String) list of Tags added to the Disk.

### Read-Only

- **id** (String) id of the Disk
- **external_id** (String) the external id of the Disk used at hypervisor

Optional:

- **create** (String)
- **delete** (String)
