package bcc_terraform

import (
	"context"
	"log"

	"fmt"
	"time"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePort() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreatePort()

	return &schema.Resource{
		CreateContext: resourcePortCreate,
		ReadContext:   resourcePortRead,
		UpdateContext: resourcePortUpdate,
		DeleteContext: resourcePortDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: args,
	}
}

func resourcePortCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}
	portNetwork, err := GetNetworkById(d, manager, nil)
	if err != nil {
		return diag.Errorf("Error getting network: %s", err)
	}

	firewallsCount := d.Get("firewall_templates.#").(int)
	firewalls := make([]*bcc.FirewallTemplate, firewallsCount)
	firewallsResourceData := d.Get("firewall_templates").(*schema.Set).List()
	for j, firewallId := range firewallsResourceData {
		portFirewall, err := manager.GetFirewallTemplate(firewallId.(string))
		if err != nil {
			return diag.Errorf("firewall_templates: Error getting Firewall Template: %s", err)
		}
		firewalls[j] = portFirewall
	}

	ipAddressInterface, ok := d.GetOk("ip_address")
	var ipAddressStr string
	if ok {
		ipAddressStr = ipAddressInterface.(string)
	} else {
		ipAddressStr = "0.0.0.0"
	}

	log.Printf("[DEBUG] subnetInfo: %#v", targetVdc)
	newPort := bcc.NewPort(portNetwork, firewalls, ipAddressStr)
	newPort.Tags = unmarshalTagNames(d.Get("tags"))
	fmt.Println(ipAddressStr)
	targetVdc.WaitLock()
	if err = targetVdc.CreateEmptyPort(&newPort); err != nil {
		return diag.Errorf("Error creating port: %s", err)
	}
	newPort.WaitLock()
	d.SetId(newPort.ID)
	fmt.Println(ipAddressStr)
	log.Printf("[INFO] Port created, ID: %s", d.Id())

	return resourcePortRead(ctx, d, meta)
}

func resourcePortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	port, err := manager.GetPort(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting port: %s", err)
		}
	}

	d.SetId(port.ID)
	d.Set("ip_address", port.IpAddress)
	d.Set("network_id", port.Network)
	d.Set("tags", marshalTagNames(port.Tags))

	firewalls := make([]*string, len(port.FirewallTemplates))
	for i, firewall := range port.FirewallTemplates {
		firewalls[i] = &firewall.ID
	}

	d.Set("firewall_templates", firewalls)

	return nil
}

func resourcePortUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	portId := d.Id()
	port, err := manager.GetPort(portId)
	if err != nil {
		return diag.Errorf("id: Error getting port: %s", err)
	}
	if d.HasChange("tags") {
		port.Tags = unmarshalTagNames(d.Get("tags"))
	}
	ip_address := d.Get("ip_address").(string)
	if d.HasChange("ip_address") {
		port.IpAddress = &ip_address
	}

	if d.HasChange("firewall_templates") {
		firewallsCount := d.Get("firewall_templates.#").(int)
		firewalls := make([]*bcc.FirewallTemplate, firewallsCount)
		firewallsResourceData := d.Get("firewall_templates").(*schema.Set).List()
		for j, firewallId := range firewallsResourceData {
			portFirewall, err := manager.GetFirewallTemplate(firewallId.(string))
			if err != nil {
				return diag.Errorf("firewall_templates: Error updating Firewall Template: %s", err)
			}
			firewalls[j] = portFirewall
		}

		port.FirewallTemplates = firewalls
	}
	if err := port.Update(); err != nil {
		return diag.FromErr(err)
	}
	port.WaitLock()
	return resourcePortRead(ctx, d, meta)
}

func resourcePortDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	portId := d.Id()

	port, err := manager.GetPort(portId)
	if err != nil {
		return diag.Errorf("id: Error getting port: %s", err)
	}

	err = port.ForceDelete()
	if err != nil {
		return diag.Errorf("Error deleting port: %s", err)
	}
	port.WaitLock()

	d.SetId("")
	log.Printf("[INFO] Port deleted, ID: %s", portId)
	return nil
}
