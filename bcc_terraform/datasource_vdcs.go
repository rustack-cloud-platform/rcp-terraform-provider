package bcc_terraform

import (
	"context"
	"fmt"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourceVdcs() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultListVdc()

	return &schema.Resource{
		ReadContext: dataSourceVdcsRead,
		Schema:      args,
	}
}

func dataSourceVdcsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetProject, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}

	allVdcs, err := manager.GetVdcs(bcc.Arguments{"project": targetProject.ID})
	if err != nil {
		return diag.Errorf("Error retrieving vdcs: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(allVdcs))
	for i, vdc := range allVdcs {
		flattenedRecords[i] = map[string]interface{}{
			"id":              vdc.ID,
			"name":            vdc.Name,
			"hypervisor":      vdc.Hypervisor.Name,
			"hypervisor_type": vdc.Hypervisor.Type,
			// "project":         vdc.Project.Name,
		}
	}

	hash, err := hashstructure.Hash(allVdcs, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `vdcs` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("vdcs/%d", hash))

	if err := d.Set("vdcs", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `vdcs` attribute: %s", err)
	}

	return nil
}
