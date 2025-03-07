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

data "basis_dns" "dns" {
    name = "test.test."
    project_id = data.basis_project.single_project.id
}

resource "basis_dns_record" "dns_record1" {
    dns_id = data.basis_dns.dns.id
    type = "A"
    host = "test.test.test."
    data = "8.8.8.8"
}