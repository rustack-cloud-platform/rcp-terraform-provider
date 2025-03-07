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

data "basis_vdc" "single_vdc" {
  project_id = data.basis_project.single_project.id
  name       = "Terraform VDC"
}

data "basis_network" "network1" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "Network 1"
}

resource "basis_port" "router_port" {
    vdc_id = resource.basis_vdc.single_vdc.id
    network_id = resource.basis_network.network1.id
}
