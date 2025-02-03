package bcc_terraform

import (
	"context"
	"log"
	"time"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceS3Storage() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectCreateS3Storage()

	return &schema.Resource{
		CreateContext: resourceS3StorageCreate,
		ReadContext:   resourceS3StorageRead,
		UpdateContext: resourceS3StorageUpdate,
		DeleteContext: resourceS3StorageDelete,
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

func resourceS3StorageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	project, err := GetProjectById(d, manager)
	if err != nil {
		return diag.Errorf("project_id: Error getting Project: %s", err)
	}
	name := d.Get("name").(string)
	backend := d.Get("backend").(string)
	newS3Storage := bcc.NewS3Storage(name, backend)
	newS3Storage.Tags = unmarshalTagNames(d.Get("tags"))

	err = project.CreateS3Storage(&newS3Storage)
	if err != nil {
		return diag.Errorf("Error creating S3Storage: %s", err)
	}

	newS3Storage.WaitLock()
	d.SetId(newS3Storage.ID)
	log.Printf("[INFO] S3Storage created, ID: %s", d.Id())

	return resourceS3StorageRead(ctx, d, meta)
}

func resourceS3StorageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	s3, err := manager.GetS3Storage(d.Id())
	if err != nil {
		return diag.Errorf("Error getting S3Storage: %s", err)
	}
	if d.HasChange("name") {
		s3.Name = d.Get("name").(string)
	}
	if d.HasChange("tags") {
		s3.Tags = unmarshalTagNames(d.Get("tags"))
	}

	err = s3.Update()
	if err != nil {
		return diag.Errorf("Error updating S3Storage: %s", err)
	}
	s3.WaitLock()
	log.Printf("[INFO] S3Storage updated, ID: %s", d.Id())

	return resourceS3StorageRead(ctx, d, meta)
}

func resourceS3StorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	S3Storage, err := manager.GetS3Storage(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting S3Storage: %s", err)
		}
	}

	d.SetId(S3Storage.ID)
	d.Set("name", S3Storage.Name)
	d.Set("backend", S3Storage.Backend)
	d.Set("project", S3Storage.Project.ID)
	d.Set("client_endpoint", S3Storage.ClientEndpoint)
	d.Set("secret_key", S3Storage.SecretKey)
	d.Set("access_key", S3Storage.AccessKey)
	d.Set("tags", marshalTagNames(S3Storage.Tags))

	return nil
}

func resourceS3StorageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	s3_id := d.Id()
	s3, err := manager.GetS3Storage(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting S3Storage: %s", err)
	}

	err = s3.Delete()
	if err != nil {
		return diag.Errorf("Error deleting S3Storage: %s", err)
	}
	s3.WaitLock()

	d.SetId("")
	log.Printf("[INFO] S3Storage deleted, ID: %s", s3_id)

	return nil
}
