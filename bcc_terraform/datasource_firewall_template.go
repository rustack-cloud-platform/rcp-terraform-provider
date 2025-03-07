package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFirewallTemplate() *schema.Resource {
	args := Defaults()
	args.injectResultFirewallTemplate()
	args.injectContextVdcById()
	args.injectContextGetFirewallTemplate() // override name

	return &schema.Resource{
		ReadContext: dataSourceFirewallTemplateRead,
		Schema:      args,
	}
}

func dataSourceFirewallTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting template: %s", err)
	}
	var targetFirewallTemplate *bcc.FirewallTemplate
	if target == "id" {
		targetFirewallTemplate, err = manager.GetFirewallTemplate(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting template: %s", err)
		}
	} else {
		targetFirewallTemplate, err = GetFirewallTemplateByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting template: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":   targetFirewallTemplate.ID,
		"name": targetFirewallTemplate.Name,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetFirewallTemplate.ID)
	return nil
}
