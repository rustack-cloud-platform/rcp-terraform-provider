---
page_title: "basis_kubernetes_templates Data Source - terraform-provider-bcc"
---
# basis_kubernetes_templates (Data Source)

Get information about kubernetes templates in the Vdc for use in other resources.

Note: You can use the [`basis_kubernetes_templates`](kubernetes_templates) data source to obtain metadata
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

data "basis_kubernetes_templates" "k8s_template" {
    vdc_id = data.basis_vdc.single_vdc.id
}

```

## Schema

### Required

- **vdc_id** (String) id of the VDC

### Read-Only

- **kubernetes_templates** (List of Object) (see [below for nested schema](#nestedatt--kubernetes_templates))

<a id="nestedatt--kubernetes_template"></a>
### Nested Schema for `kubernetes template`

Read-Only:

- **id** (String)
- **min_node_cpu** (Integer)
- **min_node_hdd** (Integer)
- **min_node_ram** (Integer)
- **name** (String)
