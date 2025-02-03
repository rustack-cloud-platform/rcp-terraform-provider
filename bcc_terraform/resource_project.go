package bcc_terraform

import (
	"context"
	"log"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProject() *schema.Resource {
	args := Defaults()
	args.injectCreateProject()

	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: args,
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	client_id := manager.ClientID
	var client *bcc.Client
	var err error

	if client_id != "" {
		client, err = manager.GetClient(client_id)
		if err != nil {
			return diag.Errorf("Error getting client: %s", err)
		}
	} else {
		allClients, err := manager.GetClients()
		if err != nil {
			return diag.Errorf("Error there are no clients available for management: %s", err)
		}
		if len(allClients) == 0 {
			return diag.Errorf("There are no available clients")
		}
		if len(allClients) > 1 {
			return diag.Errorf("More than one client available for you")
		}

		client = allClients[0]
	}

	project := bcc.NewProject(
		d.Get("name").(string),
	)
	project.Tags = unmarshalTagNames(d.Get("tags"))
	log.Printf("[DEBUG] Project create request: %#v", project)
	err = client.CreateProject(&project)
	if err != nil {
		return diag.Errorf("id: Error creating project: %s", err)
	}
	project.WaitLock()

	d.SetId(project.ID)
	log.Printf("[INFO] Project created, ID: %s", d.Id())

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	project, err := manager.GetProject(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting project: %s", err)
		}
	}

	d.SetId(project.ID)
	d.Set("name", project.Name)
	d.Set("tags", marshalTagNames(project.Tags))

	return nil
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	project, err := manager.GetProject(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting project: %s", err)
	}

	if d.HasChange("name") {
		project.Name = d.Get("name").(string)
	}
	if d.HasChange("tags") {
		project.Tags = unmarshalTagNames(d.Get("tags"))
	}
	err = project.Update()
	if err != nil {
		return diag.Errorf("name: Error rename project: %s", err)
	}
	project.WaitLock()

	log.Printf("[INFO] Updated Project, ID: %#v", project)

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	projectId := d.Id()

	project, err := manager.GetProject(projectId)
	if err != nil {
		return diag.Errorf("id: Error getting project: %s", err)
	}

	err = project.Delete()
	if err != nil {
		return diag.Errorf("Error deleting project: %s", err)
	}
	project.WaitLock()

	d.SetId("")
	log.Printf("[INFO] Project deleted, ID: %s", projectId)

	return nil
}
