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
}

data "basis_account" "me"{}

data "basis_pub_key" "key"{
    name = "name"
    # or
    or = "id"
    account_id = data.basis_account.me.id
}