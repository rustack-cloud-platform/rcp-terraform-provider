package bcc_terraform

import (
	"context"
	"log"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLbaas() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreateLbaas()

	return &schema.Resource{
		CreateContext: resourceLbaasCreate,
		ReadContext:   resourceLbaasRead,
		UpdateContext: resourceLbaasUpdate,
		DeleteContext: resourceLbaasDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceLbaasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	vdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting vdc : %s", err)
	}

	// create port
	var floatingIp *bcc.Port = nil
	if d.Get("floating").(bool) {
		floatingIp = &bcc.Port{ID: "RANDOM_FIP"}
	}
	portPrefix := "port.0"
	lbaasPort := d.Get("port.0").(map[string]interface{})

	network, err := manager.GetNetwork(lbaasPort["network_id"].(string))
	if err != nil {
		return diag.Errorf("network_id: Error getting network by id: %s", err)
	}
	network.WaitLock()
	firewalls := make([]*bcc.FirewallTemplate, 0)
	ipAddressStr := d.Get(MakePrefix(&portPrefix, "ip_address")).(string)
	if ipAddressStr == "" {
		ipAddressStr = "0.0.0.0"
	}
	port := bcc.NewPort(network, firewalls, ipAddressStr)

	newLbaas := bcc.NewLoadBalancer(d.Get("name").(string), vdc, &port, floatingIp)
	newLbaas.Tags = unmarshalTagNames(d.Get("tags"))

	err = vdc.Create(&newLbaas)
	if err != nil {
		return diag.Errorf("Error creating Lbaas: %s", err)
	}
	newLbaas.WaitLock()
	d.SetId(newLbaas.ID)
	return resourceLbaasRead(ctx, d, meta)
}

func resourceLbaasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diagErr diag.Diagnostics) {
	manager := meta.(*CombinedConfig).Manager()
	lbaas, err := manager.GetLoadBalancer(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting Lbaas: %s", err)
		}
	}
	d.SetId(lbaas.ID)
	d.Set("name", lbaas.Name)
	d.Set("floating", lbaas.Floating != nil)
	d.Set("floating_ip", "")
	if lbaas.Floating != nil {
		d.Set("floating_ip", lbaas.Floating.IpAddress)
	}
	lbaasPort := make([]interface{}, 1)
	lbaasPort[0] = map[string]interface{}{
		"ip_address": lbaas.Port.IpAddress,
		"network_id": lbaas.Port.Network.ID,
	}
	d.Set("port", lbaasPort)
	d.Set("vdc_id", lbaas.Vdc.ID)
	d.Set("tags", marshalTagNames(lbaas.Tags))

	return
}

func resourceLbaasUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	lbaas, err := manager.GetLoadBalancer(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}
	if d.HasChange("name") {
		lbaas.Name = d.Get("name").(string)
	}
	if d.HasChange("floating") {
		if !d.Get("floating").(bool) {
			lbaas.Floating = &bcc.Port{IpAddress: nil}
		} else {
			lbaas.Floating = &bcc.Port{ID: "RANDOM_FIP"}
		}
		d.Set("floating", lbaas.Floating != nil)
	}
	if d.HasChange("tags") {
		lbaas.Tags = unmarshalTagNames(d.Get("tags"))
	}
	lbaasPort := d.Get("port.0").(map[string]interface{})
	ip_address := lbaasPort["ip_address"].(string)
	if ip_address != *lbaas.Port.IpAddress {
		lbaas.Port.IpAddress = &ip_address
	}
	if err := repeatOnError(lbaas.Update, lbaas); err != nil {
		return diag.Errorf("Error updating lbaas: %s", err)
	}
	lbaas.WaitLock()

	return resourceLbaasRead(ctx, d, meta)
}

func resourceLbaasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	lbaasId := d.Id()

	lbaas, err := manager.GetLoadBalancer(lbaasId)
	if err != nil {
		return diag.Errorf("id: Error getting Lbaas: %s", err)
	}

	lbaas.Delete()
	if err != nil {
		return diag.Errorf("Error deleting Lbaas: %s", err)
	}
	lbaas.WaitLock()

	d.SetId("")
	log.Printf("[INFO] Lbaas deleted, ID: %s", lbaasId)

	return nil
}
