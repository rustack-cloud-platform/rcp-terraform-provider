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

data "basis_kubernetes_template" "kubernetes_template"{
    name = "Kubernetes 1.22.1"
    vdc_id = resource.basis_vdc.vdc1.id
    
}

data "basis_storage_profile" "ssd" {
    vdc_id = resource.basis_vdc.vdc.id
    name = "ssd"
}

data "basis_pub_key" "key"{
    name = "test"
    account_id = data.basis_account.me.id
}

data "basis_platform" "platform"{
    vdc_id = resource.basis_vdc.vdc1.id
    name = "Intel Cascade Lake"
    
}

resource "basis_kubernetes" "k8s" {
    vdc_id = resource.basis_vdc.vdc1.id
    name = "kubernetes"
    node_ram = 3
    node_cpu = 3
    platform = data.basis_platform.platform.id # vmware hypervosor only
    template_id = data.basis_kubernetes_template.kuber.id
    nodes_count = 2
    node_disk_size = 10
    node_storage_profile_id = data.basis_storage_profile.ssd.id
    user_public_key_id = data.basis_pub_key.key.id
    floating = true
}

# For get dashboard url
output "dashboard_k8s" {
    value = resource.basis_kubernetes.k8s.dashboard_url
}