package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceStorageProfile() *schema.Resource {
	args := Defaults()
	args.injectResultStorageProfile()
	args.injectContextVdcById()
	args.injectContextGetStorageProfile() // override name

	return &schema.Resource{
		ReadContext: dataSourceStorageProfileRead,
		Schema:      args,
	}
}

func dataSourceStorageProfileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting storage profile: %s", err)
	}
	var targetStorageProfile *bcc.StorageProfile
	if target == "id" {
		targetStorageProfile, err = targetVdc.GetStorageProfile(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting storage profile: %s", err)
		}
	} else {
		targetStorageProfile, err = GetStorageProfileByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting storage profile: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":   targetStorageProfile.ID,
		"name": targetStorageProfile.Name,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetStorageProfile.ID)
	return nil
}
