package bcc_terraform

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func (args *Arguments) injectContextGetProject() {
	args.merge(Arguments{
		"name": {
			Type:        schema.TypeString,
			Description: "name of the Project",
			Optional:    true,
			Computed:    true,
		},
		"id": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "id of the Project",
		},
	})
}

func (args *Arguments) injectContextProjectById() {
	args.merge(Arguments{
		"project_id": {
			Type:        schema.TypeString,
			Description: "id of the Project",
			Required:    true,
			ForceNew:    true,
		},
	})
}

func (args *Arguments) injectContextProjectByIdOptional() {
	args.merge(Arguments{
		"project_id": {
			Type:        schema.TypeString,
			Description: "id of the Project",
			Optional:    true,
		},
	})
}

func (args *Arguments) injectCreateProject() {
	args.merge(Arguments{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
				validation.StringLenBetween(1, 100),
			),
			Description: "name of the Project",
		},
		"tags": newTagNamesResourceSchema("tags of the Project"),
	})
}

func (args *Arguments) injectResultProject() {
	args.merge(Arguments{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "id of the Project",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "name of the Project",
		},
	})
}

func (args *Arguments) injectResultListProject() {
	s := Defaults()
	s.injectResultProject()

	args.merge(Arguments{
		"projects": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: s,
			},
		},
	})
}
