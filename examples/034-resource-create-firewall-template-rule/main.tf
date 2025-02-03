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

resource "basis_firewall_template" "single_template" {
  vdc_id = data.basis_vdc.single_vdc.id
  name   = "New custom template"
}

resource "basis_firewall_template_rule" "rule_1" {
    firewall_id = resource.basis_firewall_template.single_template.id
    name = "test"
    direction = "ingress"
    protocol = "tcp"
    port_range = "80"
    destination_ip = "0.0.0.0/0"
}

resource "basis_firewall_template_rule" "rule_2" {
    firewall_id = resource.basis_firewall_template.single_template.id
    name = "test"
    direction = "egress"
    protocol = "tcp"
    destination_ip = "0.0.0.0/0"
}
