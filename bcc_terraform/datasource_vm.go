package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVm() *schema.Resource {
	args := Defaults()
	args.injectResultVm()
	args.injectContextVdcById()
	args.injectContextGetVm() // override "name"

	return &schema.Resource{
		ReadContext: dataSourceVmRead,
		Schema:      args,
	}
}

func dataSourceVmRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}
	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting vm: %s", err)
	}
	var targetVm *bcc.Vm
	if target == "id" {
		targetVm, err = manager.GetVm(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting vm: %s", err)
		}
	} else {
		targetVm, err = GetVmByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting vm: %s", err)
		}
	}
	flattenPorts := make([]map[string]interface{}, 0, len(targetVm.Ports))
	for _, port := range targetVm.Ports {
		flattenPorts = append(flattenPorts, map[string]interface{}{
			"id":         port.ID,
			"ip_address": port.IpAddress,
		})
	}

	flatten := map[string]interface{}{
		"id":            targetVm.ID,
		"name":          targetVm.Name,
		"cpu":           targetVm.Cpu,
		"ram":           targetVm.Ram,
		"template_id":   targetVm.Template.ID,
		"template_name": targetVm.Template.Name,
		"power":         targetVm.Power,
		"floating":      nil,
		"floating_ip":   nil,
		"ports":         flattenPorts,
	}

	if targetVm.Floating != nil {
		flatten["floating"] = true
		flatten["floating_ip"] = targetVm.Floating.IpAddress
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetVm.ID)
	return nil
}
