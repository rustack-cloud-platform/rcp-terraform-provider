package bcc_terraform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/hashstructure/v2"
)

func dataSourcePlatforms() *schema.Resource {
	args := Defaults()
	args.injectResultListPlatforms()
	args.injectContextVdcById()

	return &schema.Resource{
		ReadContext: dataSourcePlatformsRead,
		Schema:      args,
	}
}

func dataSourcePlatformsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	platforms, err := manager.GetPlatforms(targetVdc.ID)
	if err != nil {
		return diag.Errorf("Error retrieving platforms: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(platforms))
	for i, platform := range platforms {
		flattenedRecords[i] = map[string]interface{}{
			"id":   platform.ID,
			"name": platform.Name,
		}
	}

	hash, err := hashstructure.Hash(platforms, hashstructure.FormatV2, nil)
	if err != nil {
		diag.Errorf("unable to set `platforms` attribute: %s", err)
	}

	d.SetId(fmt.Sprintf("platforms/%d", hash))

	if err := d.Set("platforms", flattenedRecords); err != nil {
		return diag.Errorf("unable to set `platforms` attribute: %s", err)
	}

	return nil
}
