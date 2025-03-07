---
page_title: "basis_firewall_template Resource - terraform-provider-bcc"
---
# basis_firewall_template (Resource)

Firewall allow you to manage your network traffic.

## Example Usage

```hcl

data "basis_project" "single_project" {
    name = "Terraform Project"
}

data "basis_vdc" "single_vdc" {
    project_id = data.basis_project.single_project.id
    name = "Terraform VDC"
}

data "basis_firewall_template" "single_template" {
  vdc_id = data.basis_vdc.single_vdc.id
  name   = "New custom template"
}

resource "basis_firewall_template_rule" "rule_1" {
    firewall_id = resource.basis_firewall_template.single_template.id
    name = "test1"
    direction = "ingress"
    protocol = "tcp"
    port_range = "80"
    destination_ip = "0.0.0.0"
}

```

## Schema

### Required

- **name** (String) name of the FirewallRule
- **firewall_id** (String) id of the firewallTemplate
- **direction** (String) direction of the FirewallRule.
   Can be chosen **ingress**, **egress**
- **protocol** (String) protocol of the FirewallRule.
   Can be chosen **tcp**, **udp**, **icmp**, **any**

> for protocols **tcp** and **udp** parameters are required to
  **port_range** (String) The range of ports can be only a single **number** and **{number}:{number}** or can be empty 
