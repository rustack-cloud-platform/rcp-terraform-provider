---
page_title: "basis_pub_key Data Source - terraform-provider-bcc"
---
# basis_pub_key (Data Source)

Get information about a public key for use in other resources. 

## Example Usage

```hcl

data "basis_account" "me"{}

data "basis_pub_key" "key" {
    account_id = data.basis_account.me.id
    
    name = "Debian 10"
    # or
    id = "id"
}

```

## Schema

### Required

- **name** (String) name of the public key `or` **id** (String) id of the public key
- **account_id** (String) id of the account

### Read-Only

- **fingerprint** (Integer) fingerprint of public key
- **public_key** (Integer) public_key value of public key data source
