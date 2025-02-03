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

data "basis_hypervisor" "vmware" {
    project_id = resource.basis_project.project1.id
    name = "VMware"
}

resource "basis_vdc" "vdc" {
    name = "Vmware Terraform"
    project_id = resource.basis_project.project1.id
    hypervisor_id = data.basis_hypervisor.vmware.id
}

data "basis_kubernetess" "kubernetes_list"{
    vdc_id = resource.basis_vdc.vdc1.id
}
