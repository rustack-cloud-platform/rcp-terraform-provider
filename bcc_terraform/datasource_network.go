package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetwork() *schema.Resource {
	args := Defaults()
	args.injectResultNetwork()
	args.injectContextVdcByIdForData()
	args.injectContextGetNetwork() // override "name"

	return &schema.Resource{
		ReadContext: dataSourceNetworkRead,
		Schema:      args,
	}
}

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("Error getting vdc: %s", err)
	}

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting network: %s", err)
	}
	var targetNetwork *bcc.Network
	if target == "id" {
		targetNetwork, err = manager.GetNetwork(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting network: %s", err)
		}
	} else {
		targetNetwork, err = GetNetworkByName(d, manager, targetVdc)
		if err != nil {
			return diag.Errorf("Error getting network: %s", err)

		}
	}

	allSubnets, err := targetNetwork.GetSubnets()
	if err != nil {
		return diag.Errorf("Error getting subnets")
	}

	flatten2 := make([]map[string]interface{}, len(allSubnets))
	for i, subnet := range allSubnets {
		dnsStrings := make([]string, len(subnet.DnsServers))
		for i3, dns := range subnet.DnsServers {
			dnsStrings[i3] = dns.DNSServer
		}

		flatten2[i] = map[string]interface{}{
			"id":       subnet.ID,
			"cidr":     subnet.CIDR,
			"dhcp":     subnet.IsDHCP,
			"gateway":  subnet.Gateway,
			"start_ip": subnet.StartIp,
			"end_ip":   subnet.EndIp,
			"dns":      dnsStrings,
		}
	}

	flatten := map[string]interface{}{
		"id":      targetNetwork.ID,
		"name":    targetNetwork.Name,
		"subnets": flatten2,
		"mtu":     targetNetwork.Mtu,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(targetNetwork.ID)
	return nil
}
