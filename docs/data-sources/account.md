---
page_title: "basis_account Data Source - terraform-provider-bcc"
---
# basis_account (Data Source)

Get information about a Acconut for use in other resources. 

## Example Usage

```hcl

data "basis_account" "account" { }

```
## Schema

### Read-Only

- **email** (String) The email address of current user
- **id** (String) The identifier for the current user
- **username** (String) The username of current user
