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

data "basis_hypervisor" "single_hypervisor" {
    project_id = data.basis_project.single_project.id
    name = "VMWARE"
}

resource "basis_vdc" "vdc1" {
    name = "Terraform VDC"
    project_id = data.basis_project.single_project.id
    hypervisor_id = data.basis_hypervisor.single_hypervisor.id
}
