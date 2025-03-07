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

data "basis_storage_profile" "single_storage_profile" {
  vdc_id = data.basis_vdc.single_vdc.id
  name   = "sas"
}

resource "basis_disk" "disk2" {
  vdc_id = data.basis_vdc.single_vdc.id

  name               = "Disk 1"
  storage_profile_id = data.basis_storage_profile.single_storage_profile.id
  size               = 1
}
