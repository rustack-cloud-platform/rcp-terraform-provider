package bcc_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceFirewallTemplates() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultListFirewallTemplate()

	return &schema.Resource{
		ReadContext: dataSourceFirewallTemplatesRead,
		Schema:      args,
	}
}

func dataSourceFirewallTemplatesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	allFirewallTemplates, err := targetVdc.GetFirewallTemplates()
	if err != nil {
		return diag.Errorf("Error retrieving firewall templates: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allFirewallTemplates))

	hash, err := hashstructure.Hash(allFirewallTemplates, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `firewall_templates` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("firewall_templates/%d", hash))

	if err := d.Set("firewall_templates", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `firewall_templates` attribute: %s", err)
	}

	return nil
}
