package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHypervisor() *schema.Resource {
	args := Defaults()
	args.injectResultHypervisor()
	args.injectContextProjectById()
	args.injectContextGetHypervisor()

	return &schema.Resource{
		ReadContext: dataSourceHypervisorRead,
		Schema:      args,
	}
}

func dataSourceHypervisorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetProject, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting hypervisor: %s", err)
	}
	var targetHypervisor *bcc.Hypervisor
	if target == "id" {
		targetHypervisor, err = GetHypervisorByIdRead(d, manager, targetProject)
		if err != nil {
			return diag.Errorf("Error getting hypervisor: %s", err)
		}
	} else {
		targetHypervisor, err = GetHypervisorByName(d, manager, targetProject)
		if err != nil {
			return diag.Errorf("Error getting hypervisor: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":         targetHypervisor.ID,
		"name":       targetHypervisor.Name,
		"type":       targetHypervisor.Type,
		"project_id": targetProject.ID,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetHypervisor.ID)
	return nil
}
