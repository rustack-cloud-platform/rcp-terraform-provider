terraform {
    required_version = ">= 1.0.0"

    required_providers {
        basis = {
            source  = "basis-cloud/bcc"
            version = "~> 0.1"
        }
    }
}

provider "basis" {
    api_endpoint = var.basis_endpoint
    token = var.basis_token
}

resource "basis_project" "project" {
    name = "Terraform GitLab Demo"
}

data "basis_hypervisor" "vmware" {
    project_id = basis_project.project.id
    name = "vmware"
}

resource "basis_vdc" "vdc" {
    name = "Gitlab"
    project_id = basis_project.project.id
    hypervisor_id = data.basis_hypervisor.vmware.id
}


data "basis_firewall_template" "allow_default" {
    vdc_id = basis_vdc.vdc.id
    name = "По-умолчанию"
}

data "basis_firewall_template" "allow_web" {
    vdc_id = basis_vdc.vdc.id
    name = "Разрешить WEB"
}

data "basis_firewall_template" "allow_ssh" {
    vdc_id = basis_vdc.vdc.id
    name = "Разрешить SSH"
}

data "basis_storage_profile" "ssd" {
    vdc_id = basis_vdc.vdc.id
    name = "ssd"
}

data "basis_network" "service_network" {
    vdc_id = basis_vdc.vdc.id
    name = "Сеть"
}

data "basis_template" "ubuntu20" {
    vdc_id = basis_vdc.vdc.id
    name = "Ubuntu 20.04"
}

resource "random_password" "password" {
    length           = 16
    special          = true
    override_special = "_-#"
}

data "template_file" "cloud_init" {
    template = file("cloud-config.tpl")
    vars = {
        user_login        = var.user_login
        public_key        = file(var.public_key)
        hostname          = "gitlab"
        gitlab_password   = random_password.password.result
    }
}

resource "basis_vm" "gitlab" {
    vdc_id = basis_vdc.vdc.id

    name = "GitLab"
    cpu = 8
    ram = 16

    template_id = data.basis_template.ubuntu20.id

    user_data = data.template_file.cloud_init.rendered

    disk {
        name = "Root"
        size = 50
        storage_profile_id = data.basis_storage_profile.ssd.id
    }

    port {
        network_id = data.basis_network.service_network.id
        firewall_templates = [
            data.basis_firewall_template.allow_default.id,
            data.basis_firewall_template.allow_web.id,
            data.basis_firewall_template.allow_ssh.id
        ]
    }

    floating = true
}

output "gitlab_ip" {
  value = basis_vm.gitlab.floating_ip
}

output "gitlab_user" {
  value = "root"
}

output "gitlab_password" {
  value = nonsensitive(random_password.password.result)
}
