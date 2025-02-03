package bcc_terraform

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewallRule() *schema.Resource {
	args := Defaults()
	args.injectContextFirewallTemplateById()
	args.injectCreateFirewallRule()

	return &schema.Resource{
		CreateContext: resourceFirewallRuleCreate,
		ReadContext:   resourceFirewallRuleRead,
		UpdateContext: resourceFirewallRuleUpdate,
		DeleteContext: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: args,
	}
}

func resourceFirewallRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	firewall_id := d.Get("firewall_id").(string)
	firewall, err := manager.GetFirewallTemplate(firewall_id)
	if err != nil {
		return diag.Errorf("firewall_id: Error getting FirewallTemplate: %s", err)
	}
	protocol := d.Get("protocol").(string)
	var newFirewallRule bcc.FirewallRule
	newFirewallRule.Name = d.Get("name").(string)
	newFirewallRule.Direction = d.Get("direction").(string)
	newFirewallRule.Protocol = d.Get("protocol").(string)
	newFirewallRule.DestinationIp = d.Get("destination_ip").(string)
	if protocol == "tcp" || protocol == "udp" {
		err = setUpRule(&newFirewallRule, d)
		if err != nil {
			return diag.Errorf("port_range: Error creating FirewallRule: %s", err)
		}
	}
	if err = firewall.CreateFirewallRule(&newFirewallRule); err != nil {
		return diag.Errorf("Error creating FirewallRule: %s", err)
	}
	d.SetId(newFirewallRule.ID)
	log.Printf("[INFO] Firewall Rule created, ID: %s", d.Id())
	return resourceFirewallRuleRead(ctx, d, meta)
}

func resourceFirewallRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	firewall_id := d.Get("firewall_id").(string)
	firewallRule_id := d.Id()

	firewall, err := manager.GetFirewallTemplate(firewall_id)
	if err != nil {
		return diag.Errorf("firewall_id: Error getting Firewall Template: %s", err)
	}

	firewallRule, err := firewall.GetRuleById(firewallRule_id)
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting fierwall Rule: %s", err)
		}
	}

	d.SetId(firewallRule.ID)
	d.Set("direction", firewallRule.Direction)
	d.Set("name", firewallRule.Name)
	d.Set("destination_ip", firewallRule.DestinationIp)
	d.Set("protocol", firewallRule.Protocol)
	if firewallRule.DstPortRangeMin != nil {
		d.Set("port_range", fmt.Sprintf("%d", *firewallRule.DstPortRangeMin))
	}
	if firewallRule.DstPortRangeMax != nil {
		d.Set("port_range", fmt.Sprintf("%s:%d", d.Get("port_range").(string), *firewallRule.DstPortRangeMax))
	}

	return nil
}

func resourceFirewallRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	firewall_id := d.Get("firewall_id").(string)
	firewallRule_id := d.Id()

	firewall, err := manager.GetFirewallTemplate(firewall_id)
	if err != nil {
		return diag.Errorf("firewall_id: Error getting Firewall Template: %s", err)
	}

	firewallRule, err := firewall.GetRuleById(firewallRule_id)
	if err != nil {
		return diag.Errorf("id: Error getting fierwall Rule: %s", err)
	}

	firewallRule.Name = d.Get("name").(string)
	protocol := d.Get("protocol").(string)
	firewallRule.DestinationIp = d.Get("destination_ip").(string)
	firewallRule.Protocol = d.Get("protocol").(string)
	if protocol == "tcp" || protocol == "udp" {
		err = setUpRule(firewallRule, d)
		if err != nil {
			return diag.Errorf("port_range: Error updating FirewallRule: %s", err)
		}
	}
	if err = firewallRule.Update(); err != nil {
		return diag.Errorf("Error updating Fierwall rule: %s", err)
	}

	return resourceFirewallRuleRead(ctx, d, meta)
}

func resourceFirewallRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	firewall_id := d.Get("firewall_id").(string)
	firewallRule_id := d.Id()

	firewall, err := manager.GetFirewallTemplate(firewall_id)
	if err != nil {
		return diag.Errorf("firewall_id: Error getting Firewall Template: %s", err)
	}

	firewallRule, err := firewall.GetRuleById(firewallRule_id)
	if err != nil {
		return diag.Errorf("id: Error getting fierwall Rule: %s", err)
	}

	err = firewallRule.Delete()
	if err != nil {
		return diag.Errorf("Error deleting Fierwall rule: %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] Fierwall rule deleted, ID: %s", firewallRule_id)
	return nil
}

func setUpRule(rule *bcc.FirewallRule, d *schema.ResourceData) (err error) {
	rule.DstPortRangeMax = nil
	rule.DstPortRangeMin = nil
	portRange := d.Get("port_range").(string)

	if portRange == "" {
		return nil
	}
	var min, max int
	var re_for_port_range = regexp.MustCompile(`(?m)^(\d+:\d+)$`)
	var re_for_port = regexp.MustCompile(`(?m)^(\d+)$`)
	if len(re_for_port_range.FindStringIndex(portRange)) > 0 {
		fmt.Sscanf(portRange, "%d:%d", &min, &max)
		rule.DstPortRangeMax = &max
		rule.DstPortRangeMin = &min
	} else if len(re_for_port.FindStringIndex(portRange)) > 0 {
		fmt.Sscanf(portRange, "%d", &min)
		rule.DstPortRangeMin = &min
	} else {
		return errors.New("PORT RANGE UNSUPPORTED FORMAT, " +
			"should be `val:val` or `val`")
	}

	return nil
}
