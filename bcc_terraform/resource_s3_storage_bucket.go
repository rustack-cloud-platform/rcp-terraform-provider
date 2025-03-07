package bcc_terraform

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceS3StorageBucket() *schema.Resource {
	args := Defaults()
	args.injectCreateS3StorageBucket()
	args.injectContextS3StorageById()

	return &schema.Resource{
		CreateContext: resourceS3StorageBucketCreate,
		ReadContext:   resourceS3StorageBucketRead,
		UpdateContext: resourceS3StorageBucketUpdate,
		DeleteContext: resourceS3StorageBucketDelete,
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

var re_for_name = regexp.MustCompile(`^[A-z0-9\-]+$`)

func resourceS3StorageBucketCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	s3_id := d.Get("s3_storage_id").(string)

	s3, err := manager.GetS3Storage(s3_id)
	if err != nil {
		return diag.Errorf("id: Error getting S3Storage: %s", err)
	}
	var S3StorageBucket bcc.S3StorageBucket
	if len(re_for_name.FindStringSubmatch(d.Get("name").(string))) > 0 {
		S3StorageBucket = bcc.NewS3StorageBucket(d.Get("name").(string))
	} else {
		return diag.Errorf("name: Wrong name format should be A-z, 1-0 and `-`")
	}

	err = s3.CreateBucket(&S3StorageBucket)
	if err != nil {
		return diag.Errorf("Error creating S3StorageBucket: %s", err)
	}

	d.SetId(S3StorageBucket.ID)
	log.Printf("[INFO] S3StorageBucket created, ID: %s", d.Id())

	return resourceS3StorageBucketRead(ctx, d, meta)
}

func resourceS3StorageBucketUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	s3_id := d.Get("s3_storage_id").(string)

	s3, err := manager.GetS3Storage(s3_id)
	if err != nil {
		return diag.Errorf("id: Error getting S3Storage: %s", err)
	}

	bucket, err := s3.GetBucket(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting S3StorageBucket: %s", err)
	}
	if d.HasChange("name") {
		if len(re_for_name.FindStringSubmatch(d.Get("name").(string))) > 0 {
			bucket.Name = d.Get("name").(string)
		} else {
			return diag.Errorf("name: Wrong name format should be A-z, 1-0 and `-`")
		}
	}

	err = bucket.Update()
	if err != nil {
		return diag.Errorf("Error updating S3StorageBucket: %s", err)
	}
	log.Printf("[INFO] S3StorageBucket updated, ID: %s", d.Id())

	return resourceS3StorageBucketRead(ctx, d, meta)
}

func resourceS3StorageBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	s3_id := d.Get("s3_storage_id").(string)

	s3, err := manager.GetS3Storage(s3_id)
	if err != nil {
		return diag.Errorf("id: Error getting S3Storage: %s", err)
	}

	bucket, err := s3.GetBucket(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting S3StorageBucket: %s", err)
		}
	}

	d.SetId(bucket.ID)
	d.Set("name", bucket.Name)
	d.Set("external_name", bucket.ExternalName)

	return nil
}

func resourceS3StorageBucketDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	s3_id := d.Get("s3_storage_id").(string)

	s3, err := manager.GetS3Storage(s3_id)
	if err != nil {
		return diag.Errorf("id: Error getting S3Storage: %s", err)
	}

	bucket, err := s3.GetBucket(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting S3StorageBucket: %s", err)
	}

	err = bucket.Delete()
	if err != nil {
		return diag.Errorf("Error deleting S3StorageBucket: %s", err)
	}

	d.SetId("")
	log.Printf("[INFO] S3StorageBucket deleted, ID: %s", s3_id)

	return nil
}
