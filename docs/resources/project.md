---
page_title: "basis_project Resource - terraform-provider-bcc"
---
# basis_project (Resource)

Projects allow you to organize your resources into groups that fit the way you work.

The Vdcs can be associated with a project:

## Example Usage

```hcl
resource "basis_project" "demo_project" {
    name = "Terraform Project"
    tags = ["created_by:terraform"]
}
```

## Schema

### Required

- **name** (String) name of the Project

### Optional

- **id** (String) The ID of this resource.
- **tags** (Toset, String) list of Tags added to the Project
