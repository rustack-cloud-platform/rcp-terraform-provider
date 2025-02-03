terraform {
  required_version = ">= 1.0.0"

  required_providers {
    basis = {
      source  = "basis-cloud/bcc"
    }
  }
}

provider "basis" {
    token = "[PLACE_YOUR_TOKEN_HERE]"
}

data "basis_project" "single_project" {
    name = "Terraform Project"
    # or
    id = "id"
}
