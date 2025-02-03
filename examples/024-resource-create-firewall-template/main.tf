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

resource "basis_firewall_template" "single_template" {
  vdc_id = data.basis_vdc.single_vdc.id
  name   = "New custom template"
}

resource "basis_firewall_template_rule" "ingress1" {
  firewall_id = resource.basis_firewall_template.single_template.id
  name = "ingress1"
  direction = "ingress"
  protocol = "tcp"
  port_range = "80"
  destination_ip = "2.0.0.0"
}

resource "basis_firewall_template_rule" "ingress2" {
  firewall_id = resource.basis_firewall_template.single_template.id
  name = "ingress2"
  direction = "ingress"
  protocol = "icmp"
  destination_ip = "1.0.0.0/24"
}

resource "basis_firewall_template_rule" "egress1" {
  firewall_id = resource.basis_firewall_template.single_template.id
  name = "egress1"
  direction = "egress"
  protocol = "udp"
  port_range = "53"
  destination_ip = "5.0.0.0/24"
}
