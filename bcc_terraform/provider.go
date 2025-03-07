package bcc_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BASIS_TOKEN", nil),
				Description: "The token key for API operations.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BASIS_API_URL", "https://cp.iteco.cloud"),
				Description: "The URL to use for the  API.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BASIS_CLIENT_ID", nil),
				Description: "The client id to use for managing instances.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"basis_account": dataSourceAccount(),

			"basis_project":              dataSourceProject(),             // 002-data-get-project +
			"basis_projects":             dataSourceProjects(),            // 003-data-get-projects +
			"basis_hypervisor":           dataSourceHypervisor(),          // 004-data-get-hypervisor +
			"basis_hypervisors":          dataSourceHypervisors(),         // 005-data-get-hypervisors +
			"basis_vdc":                  dataSourceVdc(),                 // 007-data-get-vdc +
			"basis_vdcs":                 dataSourceVdcs(),                // 008-data-get-vdcs +
			"basis_network":              dataSourceNetwork(),             // 010-data-get-network +
			"basis_networks":             dataSourceNetworks(),            // 011-data-get-networks +
			"basis_storage_profile":      dataSourceStorageProfile(),      // 012-data-get-storage-profile +
			"basis_storage_profiles":     dataSourceStorageProfiles(),     // 013-data-get-storage-profiles +
			"basis_disk":                 dataSourceDisk(),                // 015-data-get-disk +
			"basis_disks":                dataSourceDisks(),               // 016-data-get-disks +
			"basis_template":             dataSourceTemplate(),            // 017-data-get-template +
			"basis_templates":            dataSourceTemplates(),           // 018-data-get-templates +
			"basis_firewall_template":    dataSourceFirewallTemplate(),    // 019-data-get-template +
			"basis_firewall_templates":   dataSourceFirewallTemplates(),   // 020-data-get-templates +
			"basis_vm":                   dataSourceVm(),                  // 022-data-get-vm
			"basis_vms":                  dataSourceVms(),                 // 023-data-get-vms
			"basis_router":               dataSourceRouter(),              // 025-data-get-router +
			"basis_routers":              dataSourceRouters(),             // 026-data-get-routers +
			"basis_port":                 dataSourcePort(),                // 027-data-get-port +
			"basis_ports":                dataSourcePorts(),               // 027-data-get-ports +
			"basis_dns":                  dataSourceDns(),                 // 028-data-get-dns +
			"basis_dnss":                 dataSourceDnss(),                // 028-data-get-dnss +
			"basis_lbaas":                dataSourceLbaas(),               // 028-data-get-lbaas +
			"basis_lbaass":               dataSourceLoadBalancers(),       // 028-data-get-lbaass +
			"basis_s3":                   dataSourceLoadBalancers(),       // 028-data-get-lbaass +
			"basis_s3_storage":           dataSourceS3Storage(),           // 028-data-get-s3-storage +
			"basis_s3_storages":          dataSourceS3Storages(),          // 028-data-get-s3-storages +
			"basis_kubernetes":           dataSourceKubernetes(),          // 030-resource-get-kubernetes +
			"basis_kubernetess":          dataSourceKubernetess(),         // 030-resource-get-kubernetess +
			"basis_kubernetes_template":  dataSourceKubernetesTemplate(),  // 030-resource-get-kubernetes_template +
			"basis_kubernetes_templates": dataSourceKubernetesTemplates(), // 030-resource-get-kubernetes_templates +
			"basis_pub_key":              dataSourcePublicKey(),           // 030-resource-get-pub-key +
			"basis_platform":             dataSourcePlatform(),            // 030-resource-get-platform +
			"basis_platforms":            dataSourcePlatforms(),           // 030-resource-get-platforms +
			"basis_paas_template":        dataSourcePaasTemplate(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"basis_project":                resourceProject(),          // 001-resource-create-project +
			"basis_vdc":                    resourceVdc(),              // 006-resource-create-vdc +
			"basis_network":                resourceNetwork(),          // 009-resource-create-network +
			"basis_disk":                   resourceDisk(),             // 014-resource-create-disk +
			"basis_vm":                     resourceVm(),               // 021-resource-create-vm +
			"basis_firewall_template":      resourceFirewallTemplate(), // 024-resource-create-firewall-template +
			"basis_router":                 resourceRouter(),           // 027-resource-create-router +
			"basis_port":                   resourcePort(),             // 027-resource-create-port +
			"basis_dns":                    resourceDns(),              // 028-resource-create-dns +
			"basis_dns_record":             resourceDnsRecord(),        // 028-resource-create-dns-record +
			"basis_firewall_template_rule": resourceFirewallRule(),     // 029-resource-create-firewall-rule +
			"basis_lbaas":                  resourceLbaas(),            // 029-resource-create-lbaas +
			"basis_lbaas_pool":             resourceLbaasPool(),        // 029-resource-create-lbaas-pool +
			"basis_s3_storage":             resourceS3Storage(),        // 029-resource-create-s3-storage +
			"basis_s3_storage_bucket":      resourceS3StorageBucket(),  // 029-resource-create-s3-storage-bucket +
			"basis_kubernetes":             resourceKubernetes(),       // 030-resource-create-basis-kubernetes +
			"basis_paas_service":           resourcePaasService(),
		},
	}

	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "1.6"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	config := Config{
		Token:            d.Get("token").(string),
		APIEndpoint:      d.Get("api_endpoint").(string),
		ClientID:         d.Get("client_id").(string),
		TerraformVersion: terraformVersion,
	}

	return config.Client()
}
