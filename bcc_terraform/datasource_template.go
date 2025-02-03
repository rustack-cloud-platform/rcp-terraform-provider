package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTemplate() *schema.Resource {
	args := Defaults()
	args.injectResultTemplate()
	args.injectContextVdcById()
	args.injectContextGetTemplate() // override name

	return &schema.Resource{
		ReadContext: dataSourceTemplateRead,
		Schema:      args,
	}
}

func dataSourceTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting template: %s", err)
	}
	var targetTemplate *bcc.Template
	if target == "id" {
		targetTemplate, err = manager.GetTemplate(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting template: %s", err)
		}
	} else {
		targetTemplate, err = GetTemplateByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting template: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":       targetTemplate.ID,
		"name":     targetTemplate.Name,
		"min_cpu":  targetTemplate.MinCpu,
		"min_ram":  targetTemplate.MinRam,
		"min_disk": targetTemplate.MinHdd,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetTemplate.ID)
	return nil
}
