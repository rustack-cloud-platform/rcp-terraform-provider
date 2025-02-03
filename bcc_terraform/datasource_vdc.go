package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVdc() *schema.Resource {
	args := Defaults()
	args.injectResultVdc()
	args.injectContextProjectByIdOptional()
	args.injectContextGetVdc() // override name

	return &schema.Resource{
		ReadContext: dataSourceVdcRead,
		Schema:      args,
	}
}

func dataSourceVdcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	var targetProject *bcc.Project

	if _, exists := d.GetOk("project_id"); exists {
		project, err := GetProjectById(d, manager)
		if err != nil {
			return diag.Errorf("Error getting project: %s", err)
		}

		targetProject = project
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting VDC: %s", err)
	}
	var targetVdc *bcc.Vdc
	if target == "id" {
		targetVdc, err = manager.GetVdc(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting VDC: %s", err)
		}
	} else {
		targetVdc, err = GetVdcByName(d, manager, targetProject)
		if err != nil {
			return diag.Errorf("Error getting VDC: %s", err)
		}
	}

	flattenedVdc := map[string]interface{}{
		"id":              targetVdc.ID,
		"name":            targetVdc.Name,
		"hypervisor":      targetVdc.Hypervisor.Name,
		"hypervisor_type": targetVdc.Hypervisor.Type,
	}

	if err := setResourceDataFromMap(d, flattenedVdc); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetVdc.ID)
	return nil
}
