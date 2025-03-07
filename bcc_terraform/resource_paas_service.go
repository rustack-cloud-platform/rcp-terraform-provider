package bcc_terraform

import (
	"context"
	"encoding/json"
	"time"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePaasService() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourcePaasServiceRead,
		CreateContext: resourcePaasServiceCreate,
		UpdateContext: resourcePaasServiceUpdate,
		DeleteContext: resourcePaasServiceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "id of Paas Service at BCC",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of Paas Service at BCC",
			},
			"vdc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "id of Vdc",
			},
			"paas_service_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "id of Paas Template",
			},
			"paas_service_inputs": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "inputs of Paas Service as JSON object",
			},
		},
	}
}

func ensureLocationCreated(vdcId string, manager *bcc.Manager) error {
	vdc, err := manager.GetVdc(vdcId)
	if err != nil {
		return err
	}
	if vdc.Paas != nil {
		return nil
	}
	err = manager.CreatePaasLocation(vdcId)
	if err != nil {
		return err
	}
	for {
		vdc, err := manager.GetVdc(vdcId)
		if err != nil {
			return err
		}
		if vdc.Paas != nil && !vdc.Paas.Locked {
			return nil
		}
		time.Sleep(time.Second)
	}
}

func resourcePaasServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).Manager()
	manager = manager.WithContext(ctx)
	err := ensureLocationCreated(d.Get("vdc_id").(string), manager)
	if err != nil {
		return diag.FromErr(err)
	}
	service, err := manager.GetPaasService(d.Get("id").(string))
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("Error getting Paas Service: %s", err)
		}
	}
	d.Set("id", service.ID)
	d.Set("name", service.Name)
	d.Set("vdc_id", service.Vdc.ID)
	d.Set("paas_service_id", service.PaasServiceID)
	inputsString, err := json.Marshal(service.Inputs)
	if err != nil {
		return diag.Errorf("Error marshalling Paas Service inputs: %s", err)
	}
	d.Set("paas_service_inputs", inputsString)
	d.SetId(service.ID)
	return nil
}

func resourcePaasServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	manager = manager.WithContext(ctx)
	err := ensureLocationCreated(d.Get("vdc_id").(string), manager)
	if err != nil {
		return diag.FromErr(err)
	}
	var inputs map[string]interface{}
	inputsString := d.Get("paas_service_inputs").(string)
	if err := json.Unmarshal([]byte(inputsString), &inputs); err != nil {
		return diag.Errorf("Error parsing Paas Service inputs: %s", err)
	}
	service := &bcc.PaasService{
		Name: d.Get("name").(string),
		Vdc: struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}{
			ID: d.Get("vdc_id").(string),
		},
		PaasServiceID: d.Get("paas_service_id").(int),
		Inputs:        inputs,
	}
	if err := manager.CreatePaasService(service); err != nil {
		return diag.Errorf("Error creating Paas Service: %s", err)
	}
	service.WaitLock()
	d.Set("id", service.ID)
	return resourcePaasServiceRead(ctx, d, meta)
}

func resourcePaasServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	manager = manager.WithContext(ctx)
	err := ensureLocationCreated(d.Get("vdc_id").(string), manager)
	if err != nil {
		return diag.FromErr(err)
	}
	service, err := manager.GetPaasService(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting paas service: %s", err)
	}
	var inputs map[string]interface{}
	inputsString := d.Get("paas_service_inputs").(string)
	if err := json.Unmarshal([]byte(inputsString), &inputs); err != nil {
		return diag.Errorf("Error parsing Paas Service inputs: %s", err)
	}
	service.Inputs = inputs	
	service.Name = d.Get("name").(string)
	if err := service.Update(); err != nil {
		return diag.Errorf("Error updating Paas Service: %s", err)
	}
	service.WaitLock()

	return resourcePaasServiceRead(ctx, d, meta)
}

func resourcePaasServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	manager = manager.WithContext(ctx)
	err := ensureLocationCreated(d.Get("vdc_id").(string), manager)
	if err != nil {
		return diag.FromErr(err)
	}
	service, err := manager.GetPaasService(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = manager.DeletePaasService(d.Get("id").(string))
	if err != nil {
		return diag.Errorf("Error deleting Paas Service: %s", err)
	}
	service.WaitLock()
	return nil
}
