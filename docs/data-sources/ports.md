---
page_title: "basis_ports Data Source - terraform-provider-bcc"
---
# basis_ports (Data Source)

Get information about list of Ports in the Vdc for use in other resources.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_port" "all_port" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **ports** (List of Object) (see [below for nested schema](#nestedatt--ports))

<a id="nestedatt--ports"></a>
### Nested Schema for `ports`

Read-Only:

- **id** (String)
- **network_id** (String)
- **ip_address** (String)
- **firewall_templates** (List of String)
