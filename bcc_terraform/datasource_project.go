package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProject() *schema.Resource {
	args := Defaults()
	args.injectResultProject()
	args.injectContextGetProject()

	return &schema.Resource{
		ReadContext: dataSourceProjectRead,
		Schema:      args,
	}
}

func dataSourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting project: %s", err)
	}
	var targetProject *bcc.Project
	if target == "id" {
		targetProject, err = manager.GetProject(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting project: %s", err)
		}
	} else {
		targetProject, err = GetProjectByName(d, manager)
		if err != nil {
			return diag.Errorf("Error getting project: %s", err)
		}
	}

	flattenedProject := map[string]interface{}{
		"id":   targetProject.ID,
		"name": targetProject.Name,
		// "project_id":   nil,
		// "project_name": nil,
	}

	if err := setResourceDataFromMap(d, flattenedProject); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetProject.ID)
	return nil
}
