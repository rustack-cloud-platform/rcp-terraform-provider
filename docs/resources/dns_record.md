---
page_title: "basis_dns_record Resource - terraform-provider-bcc"
---
# basis_dns_record (Resource)

Provides a Basis DNS record resource.

## Example Usage

```hcl
data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_dns" "dns" {
    name="dns.teraform."
    project_id = data.basis_project.single_project.id
}

resource "basis_dns_record" "dns_record1" {
    dns_id = data.basis_dns.dns.id
    type = "A"
    host = "test2.dns.teraform."
    data = "8.8.8.8"
}

```

## Schema

### Required

> required for all types

- **dns_id** (String) name of the Dns
- **type** (String) type of Dns record
- **host** (String) host of Dns record
- **data** (String) data of Dns record

> for type CAA parameters are required to

- **tag** (String) tag of Dns record
- **flag** (String) flag of Dns record. Can be chosen
    **0 (not critical)**, **128 (critical)**

> for type MX parameters are required to

- **Priority** (String) Priority of Dns record

> for type SRV parameters are required to

- **Priority** (String) Priority of Dns record
- **Weight** (String) Weight of Dns record
- **Port** (String) Port of Dns record

### Optional

- **id** (String) The ID of this resource.
- **timeouts** (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- **create** (String)
- **delete** (String)
