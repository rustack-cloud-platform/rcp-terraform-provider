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

data "basis_hypervisor" "vmware" {
    project_id = resource.basis_project.project1.id
    name = "VMware"
}

resource "basis_vdc" "vdc" {
    name = "Vmware Terraform"
    project_id = resource.basis_project.project1.id
    hypervisor_id = data.basis_hypervisor.vmware.id
}

data "basis_kubernetes_template" "kuber"{
    name = "Kubernetes 1.22.1"
    # or
    or = "id"
    vdc_id = resource.basis_vdc.vdc1.id
    
}