---
page_title: "basis_firewall_templates Data Source - terraform-provider-bcc"
---
# basis_firewall_templates (Data Source)

Get information about Firewall Templates in the Vdc for use in other resources.

Note: You can use the [`basis_firewall_template`](Firewall Template) data source to obtain metadata
about a single Firewall Template if you already know the `name` and `vdc_id` to retrieve.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_firewall_templates" "all_templates" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```
## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **firewall_templates** (List of Object) (see [below for nested schema](#nestedatt--templates))

<a id="nestedatt--templates"></a>
### Nested Schema for `templates`

Read-Only:

- **id** (String)
- **name** (String)
