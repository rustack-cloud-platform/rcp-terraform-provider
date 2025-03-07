package bcc_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextGetVm() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "name of the Vm",
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the vm",
		},
	})
}

func (args *Arguments) injectContextVmById() {
	args.merge(Arguments{
		"vm_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "id of the Vm",
		},
	})
}

func (args *Arguments) injectCreateVm() {
	systemDisk := Defaults()
	systemDisk.injectSystemDisk()

	args.merge(Arguments{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 100),
			),
			Description: "name of the Vm",
		},
		"cpu": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 128),
			Description:  "the number of virtual cpus",
		},
		"ram": {
			Type:        schema.TypeFloat,
			Required:    true,
			Description: "memory of the Vm in gigabytes",
		},
		"template_id": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "id of the Template",
		},
		"user_data": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "script for cloud-init",
		},
		"system_disk": {
			Type:     schema.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: systemDisk,
			},
			Description: "System disk.",
		},
		"disks": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "list of Disks attached to the Vm",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"ports": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MinItems:    1,
			MaxItems:    10,
			Description: "List of Ports connected to the Vm",
			Deprecated:  "Use networks instead of ports",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"networks": {
			Type:         schema.TypeList,
			Optional:     true,
			Computed:     true,
			ExactlyOneOf: []string{"ports", "networks"},
			MinItems:     1,
			MaxItems:     10,
			Description:  "List of Ports connected to the Vm",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Id of the Port",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "IP of the Port",
					},
				},
			},
		},
		"floating": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "enable floating ip for the Vm",
		},
		"floating_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "floating ip for the Vm. May be omitted",
		},
		"tags": newTagNamesResourceSchema("tags of the Vm"),
		"power": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "power of vw on/off",
		},
	})
}

func (args *Arguments) injectResultVm() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Vm",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Vm",
		},
		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "the number of virtual cpus",
		},
		"ram": {
			Type:        schema.TypeFloat,
			Computed:    true,
			Description: "memory of the Vm in gigabytes",
		},
		"template_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Template",
		},
		"template_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Template",
		},
		"floating": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "enable floating ip for the Vm",
		},
		"floating_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "floating_ip of the Vm. May be omitted",
		},
		"power": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "power of vw on/off",
		},
		"ports": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "List of Ports connected to the Vm",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "Id of the Port",
					},
					"ip_address": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "IP of the Port",
					},
				},
			},
		},
	})
}

func (args *Arguments) injectResultListVm() {
	s := Defaults()
	s.injectResultVm()

	args.merge(Arguments{
		"vms": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}

func (args *Arguments) injectSystemDisk() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the System Disk",
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"storage_profile_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"external_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "external id of the volume. It can be empty",
		},
	})
}
