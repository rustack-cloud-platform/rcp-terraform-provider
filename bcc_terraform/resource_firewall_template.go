package bcc_terraform

import (
	"context"
	"log"
	"time"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFirewallTemplate() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreateFirewallTemplate()

	return &schema.Resource{
		CreateContext: resourceFirewallTemplateCreate,
		ReadContext:   resourceFirewallTemplateRead,
		UpdateContext: resourceFirewallTemplateUpdate,
		DeleteContext: resourceFirewallTemplateDelete,
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

func resourceFirewallTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	newFirewallTemplate := bcc.NewFirewallTemplate(d.Get("name").(string))
	newFirewallTemplate.Tags = unmarshalTagNames(d.Get("tags"))
	err = targetVdc.CreateFirewallTemplate(&newFirewallTemplate)
	if err != nil {
		return diag.Errorf("Error creating Firewall Template: %s", err)
	}

	d.SetId(newFirewallTemplate.ID)
	log.Printf("[INFO] FirewallTemplate created, ID: %s", d.Id())

	return resourceFirewallTemplateRead(ctx, d, meta)
}

func resourceFirewallTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	firewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting Firewall Template: %s", err)
		}
	}

	d.SetId(firewallTemplate.ID)
	d.Set("name", firewallTemplate.Name)
	d.Set("tags", marshalTagNames(firewallTemplate.Tags))

	return nil
}

func resourceFirewallTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	firewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting FirewallTemplate: %s", err)
	}

	if d.HasChange("name") {
		firewallTemplate.Name = d.Get("name").(string)
	}
	if d.HasChange("tags") {
		firewallTemplate.Tags = unmarshalTagNames(d.Get("tags"))
	}
	if err = firewallTemplate.UpdateFirewallTemplate(); err != nil {
		return diag.Errorf("name: Error rename Firewall Template: %s", err)
	}

	return resourceFirewallTemplateRead(ctx, d, meta)
}

func resourceFirewallTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	FirewallTemplate, err := manager.GetFirewallTemplate(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting FirewallTemplate: %s", err)
	}

	err = FirewallTemplate.Delete()
	if err != nil {
		return diag.Errorf("Error deleting FirewallTemplate: %s", err)
	}

	return nil
}
