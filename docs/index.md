---
page_title: "BCC provider"
---
# BCC provider

The BCC provider is used to interact with the Basis cloud. 
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
terraform {
  required_providers {
    basis = {
      source  = "basis-cloud/bcc"
    }
  }
}

# Set the variable value in *.tfvars file
# or using -var="basis_token=..." CLI option
variable "basis_token" {}

# Configure the BCC provider
provider "basis" {
    api_endpoint = "https://cp.iteco.cloud"
    token = var.basis_token
}

```

-> **Note for Module Developers** Although provider configurations are shared between modules, each module must
declare its own [provider requirements](https://www.terraform.io/docs/language/providers/requirements.html). See the [module development documentation](https://www.terraform.io/docs/language/modules/develop/providers.html) for additional information.

## Schema

### Optional

- **api_endpoint** (String) The URL to use for the BCC API.
- **token** (String) The token key for API operations.
