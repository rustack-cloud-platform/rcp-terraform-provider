package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePort() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectResultPort()

	return &schema.Resource{
		ReadContext: dataSourcePortRead,
		Schema:      args,
	}
}

func dataSourcePortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	port_id := d.Get("id")
	portIp := d.Get("ip_address")
	flatten := make(map[string]interface{})
	var targetPort *bcc.Port

	if port_id != "" && portIp != "" {
		return diag.Errorf("For getting the port must be specified id or ip")
	}

	if port_id != "" {
		targetPort, err = GetPortById(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting port: %s", err)
		}
	} else {
		targetPort, err = GetPortByIp(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting port: %s", err)
		}

	}

	flatten["id"] = targetPort.ID
	flatten["ip_address"] = targetPort.IpAddress
	flatten["network"] = targetPort.Network.ID

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetPort.ID)
	return nil
}
