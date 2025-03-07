package bcc_terraform

import (
	"context"
	"fmt"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePaasTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePaasTemplateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "id of Paas Template",
			},
			"vdc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "id of Vdc",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "name of Paas Template",
				AtLeastOneOf: []string{"id", "name"},
			},
		},
	}
}

func dataSourcePaasTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	manager = manager.WithContext(ctx)
	err := ensureLocationCreated(d.Get("vdc_id").(string), manager)
	if err != nil {
		return diag.FromErr(err)
	}
	idValue, ok := d.GetOk("id")
	var template *bcc.PaasTemplate
	if ok {
		template, err = manager.GetPaasTemplate(idValue.(int), d.Get("vdc_id").(string))
		if err != nil {
			if err.(*bcc.ApiError).Code() == 404 {
				d.SetId("")
				return nil
			}
			return diag.Errorf("Error getting paas template: %s", err)
		}
	} else {
		templates, err := manager.GetPaasTemplates(d.Get("vdc_id").(string), bcc.Arguments{"name": d.Get("name").(string)})
		if err != nil {
			return diag.Errorf("Error getting paas template: %s", err)
		}
		if len(templates) == 0 {
			d.SetId("")
			return nil
		}
		template = templates[0]
	}
	flatten := map[string]interface{}{
		"id":   template.ID,
		"name": template.Name,
	}
	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprint(template.ID))
	return nil
}
