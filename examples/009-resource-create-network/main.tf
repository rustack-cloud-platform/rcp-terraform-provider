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
    name = "Terraform VDC"
}

resource "basis_network" "network2" {
    vdc_id = data.basis_vdc.single_vdc.id

    name = "Сеть 1"

    subnets {
        cidr = "10.20.40.0/24"
        dhcp = true
        gateway = "10.20.40.1"
        start_ip = "10.20.40.2"
        end_ip = "10.20.40.254"
        dns = ["8.8.8.8", "8.8.4.4", "1.1.1.1"]
    }
}
