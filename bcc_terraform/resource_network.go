package bcc_terraform

import (
	"context"
	"log"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetwork() *schema.Resource {
	args := Defaults()
	args.injectContextVdcById()
	args.injectCreateNetwork()

	return &schema.Resource{
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		UpdateContext: resourceNetworkUpdate,
		DeleteContext: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	log.Printf("[DEBUG] subnetInfo: %#v", targetVdc)
	network := bcc.NewNetwork(d.Get("name").(string))
	network.Tags = unmarshalTagNames(d.Get("tags"))
	if mtu, ok := d.GetOk("mtu"); ok {
		mtuValue := mtu.(int)
		network.Mtu = &mtuValue
	} else {
		network.Mtu = nil
	}
	targetVdc.WaitLock()
	if err = targetVdc.CreateNetwork(&network); err != nil {
		return diag.Errorf("Error creating network: %s", err)
	}
	d.SetId(network.ID)

	diag := createSubnet(d, manager)
	if diag != nil {
		return diag
	}
	network.WaitLock()

	log.Printf("[INFO] Network created, ID: %s", d.Id())

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting network: %s", err)
		}
	}

	d.Set("name", network.Name)
	d.Set("tags", marshalTagNames(network.Tags))
	d.Set("mtu", network.Mtu)

	subnets, err := network.GetSubnets()
	if err != nil {
		return diag.Errorf("subnets: Error getting subnets: %s", err)
	}

	flattenedRecords := make([]map[string]interface{}, len(subnets))
	for i, subnet := range subnets {
		dnsStrings := make([]string, len(subnet.DnsServers))
		for i2, dns := range subnet.DnsServers {
			dnsStrings[i2] = dns.DNSServer
		}
		flattenedRecords[i] = map[string]interface{}{
			"id":       subnet.ID,
			"cidr":     subnet.CIDR,
			"dhcp":     subnet.IsDHCP,
			"gateway":  subnet.Gateway,
			"start_ip": subnet.StartIp,
			"end_ip":   subnet.EndIp,
			"dns":      dnsStrings,
		}
	}

	if err := d.Set("subnets", flattenedRecords); err != nil {
		return diag.Errorf("subnets: unable to set `subnet` attribute: %s", err)
	}

	return nil
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting network: %s", err)
	}
	shouldUpdate := false
	if d.HasChange("tags") {
		network.Tags = unmarshalTagNames(d.Get("tags"))
		shouldUpdate = true
	}

	if d.HasChange("name") {
		network.Name = d.Get("name").(string)
		shouldUpdate = true
	}
	if d.HasChange("mtu") {
		if mtu, ok := d.GetOk("mtu"); ok {
			mtuValue := mtu.(int)
			network.Mtu = &mtuValue
		} else {
			network.Mtu = nil
		}
		shouldUpdate = true
	}
	if shouldUpdate {
		err := network.Update()
		if err != nil {
			return diag.Errorf("name: Error update network: %s", err)
		}
	}

	if d.HasChange("subnets") {
		diagErr := updateSubnet(d, manager)
		if diagErr != nil {
			return diagErr
		}
	}
	network.WaitLock()

	return resourceNetworkRead(ctx, d, meta)
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting network: %s", err)

	}

	if err = repeatOnError(network.Delete, network); err != nil {
		return diag.Errorf("Error deleting network: %s", err)
	}
	network.WaitLock()

	return nil
}

func createSubnet(d *schema.ResourceData, manager *bcc.Manager) (diagErr diag.Diagnostics) {
	subnets := d.Get("subnets").([]interface{})
	log.Printf("[DEBUG] subnets: %#v", subnets)
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Unable to get network: %s", err)
	}

	for _, subnetInfo := range subnets {
		log.Printf("[DEBUG] subnetInfo: %#v", subnetInfo)
		subnetInfo2 := subnetInfo.(map[string]interface{})

		// Create subnet
		subnet := bcc.NewSubnet(subnetInfo2["cidr"].(string), subnetInfo2["gateway"].(string), subnetInfo2["start_ip"].(string), subnetInfo2["end_ip"].(string), subnetInfo2["dhcp"].(bool))

		if err := network.CreateSubnet(&subnet); err != nil {
			return diag.Errorf("subnets: Error creating subnet: %s", err)
		}

		dnsServersList := subnetInfo2["dns"].([]interface{})
		dnsServers := make([]*bcc.SubnetDNSServer, len(dnsServersList))
		for i, dns := range dnsServersList {
			s1 := bcc.NewSubnetDNSServer(dns.(string))
			dnsServers[i] = &s1
		}

		if err := subnet.UpdateDNSServers(dnsServers); err != nil {
			return diag.Errorf("dns: Error Update DNS Servers: %s", err)
		}

	}

	return
}

func updateSubnet(d *schema.ResourceData, manager *bcc.Manager) (diagErr diag.Diagnostics) {

	subnets := d.Get("subnets").([]interface{})
	network, err := manager.GetNetwork(d.Id())
	if err != nil {
		return diag.Errorf("id: Unable to get network: %s", err)
	}
	subnetsRaw, err := network.GetSubnets()
	if err != nil {
		return diag.Errorf("subnets: Unable to get subnets: %s", err)
	}

	for _, subnetInfo := range subnets {
		subnetInfo2 := subnetInfo.(map[string]interface{})
		var subnet *bcc.Subnet
		for _, currentSubnet := range subnetsRaw {
			if currentSubnet.CIDR == subnetInfo2["cidr"] {
				subnet = currentSubnet
				break
			}
		}
		dnsServersList := subnetInfo2["dns"].([]interface{})
		newDnsServers := make([]*bcc.SubnetDNSServer, len(dnsServersList))
		for i, dns := range dnsServersList {
			s1 := bcc.NewSubnetDNSServer(dns.(string))
			newDnsServers[i] = &s1
		}
		if subnet == nil {
			// create new subnet
			newSubnet := bcc.NewSubnet(subnetInfo2["cidr"].(string), subnetInfo2["gateway"].(string), subnetInfo2["start_ip"].(string), subnetInfo2["end_ip"].(string), subnetInfo2["dhcp"].(bool))
			if err := network.CreateSubnet(&newSubnet); err != nil {
				return diag.Errorf("subnets: Error creating subnet: %s", err)
			}
			if err := subnet.UpdateDNSServers(newDnsServers); err != nil {
				return diag.Errorf("dns: Error Update DNS Servers: %s", err)
			}
		} else {
			// update preserved subnet
			shouldUpdate := false
			if subnet.Gateway != subnetInfo2["gateway"] {
				return diag.Errorf("You cannot change gateway")
			}
			if subnet.StartIp != subnetInfo2["start_ip"] || subnet.EndIp != subnetInfo2["end_ip"] || subnet.IsDHCP != subnetInfo2["dhcp"] {
				subnet.EndIp = subnetInfo2["end_ip"].(string)
				subnet.StartIp = subnetInfo2["start_ip"].(string)
				subnet.IsDHCP = subnetInfo2["dhcp"].(bool)
				shouldUpdate = true
			}
			if len(subnet.DnsServers) != len(newDnsServers) {
				subnet.DnsServers = newDnsServers
				shouldUpdate = true
			} else {
				for i, oldDns := range subnet.DnsServers {
					if oldDns.DNSServer != newDnsServers[i].DNSServer {
						subnet.DnsServers = newDnsServers
						shouldUpdate = true
						break
					}
				}
			}
			if shouldUpdate {
				if err := subnet.UpdateDNSServers(subnet.DnsServers); err != nil {
					return diag.Errorf("error update subnet: %s", err)
				}
			}
		}
	}
	for _, subnet := range subnetsRaw {
		var subnetInfo2 map[string]interface{}
		for _, subnetInfo := range subnets {
			subnetInfo2 = subnetInfo.(map[string]interface{})
			if subnet.CIDR == subnetInfo2["cidr"] {
				break
			}
		}
		if subnetInfo2 == nil {
			// delete obsolete subnet
			if err := subnet.Delete(); err != nil {
				return diag.Errorf("error deleting subnet: %s", err)
			}
		}
	}

	return
}
