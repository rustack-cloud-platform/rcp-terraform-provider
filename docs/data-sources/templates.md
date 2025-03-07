---
page_title: "basis_templates Data Source - terraform-provider-bcc"
---
# basis_templates (Data Source)

Get information about Templates in the Vdc for use in other resources.

Note: You can use the [`basis_template`](template) data source to obtain metadata
about a single Template if you already know the `name` and `vdc_id` to retrieve.


## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_templates" "single_template" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **templates** (List of Object) (see [below for nested schema](#nestedatt--templates))

<a id="nestedatt--templates"></a>
### Nested Schema for `templates`

Read-Only:

- **id** (String)
- **min_cpu** (Integer)
- **min_disk** (Integer)
- **min_ram** (Integer)
- **name** (String)
