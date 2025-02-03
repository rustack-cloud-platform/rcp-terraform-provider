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

data "basis_network" "service_network" {
    vdc_id = data.basis_vdc.single_vdc.id
    name = "Сеть"
}

data "basis_firewall_template" "allow_default" {
    vdc_id = resource.basis_vdc.vdc1.id
    name = "Разрешить входящие"
}

data "basis_storage_profile" "ssd" {
    vdc_id = resource.basis_vdc.vdc1.id
    name = "ssd"
}


resource "basis_port" "vm_port" {
    vdc_id = resource.basis_vdc.vdc1.id

    network_id = resource.basis_network.service_network.id
    firewall_templates = [data.basis_firewall_template.allow_default.id]
}

resource "basis_vm" "vm" {
    vdc_id = resource.basis_vdc.vdc1.id
    name = "Server 1"
    cpu = 3
    ram = 3
    power = true

    template_id = data.basis_template.ubuntu20.id

    user_data = "${file("user_data.yaml")}"

    system_disk {
        size = 10
        storage_profile_id = data.basis_storage_profile.ssd.id
    }

    ports = [resource.basis_port.vm_port.id]

    floating = true
}


resource "basis_lbaas" "lbaas" {
    vdc_id = resource.basis_vdc.single_vdc.id
    name = "lbaas"
    floating = true
    port {
        network_id = data.basis_network.service_network.id
    }
}

resource "basis_lbaas_pool" "pool" {
    lbaas_id = resource.basis_lbaas.lbaas.id
    connlimit = 34
    method = "SOURCE_IP"
    protocol = "TCP"
    port = 80

     member {
        port = 80
        weight = 50
        vm_id = resource.basis_vm.vm.id
    }

    depends_on = [basis_vm.vm1]
}
