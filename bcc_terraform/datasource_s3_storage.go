package bcc_terraform

import (
	"context"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceS3Storage() *schema.Resource {
	args := Defaults()
	args.injectContextProjectById()
	args.injectResultS3Storage()
	args.injectContextGetS3Storage() // override name

	return &schema.Resource{
		ReadContext: dataSourceS3StorageRead,
		Schema:      args,
	}
}

func dataSourceS3StorageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()

	target, err := checkDatasourceNameOrId(d)
	if err != nil {
		return diag.Errorf("Error getting s3 storage: %s", err)
	}
	var s3_storage *bcc.S3Storage
	if target == "id" {
		s3_storage, err = manager.GetS3Storage(d.Get("id").(string))
		if err != nil {
			return diag.Errorf("Error getting storage: %s", err)
		}
	} else {
		s3_storage, err = GetS3ByName(d, manager)
		if err != nil {
			return diag.Errorf("Error getting storage: %s", err)
		}
	}

	flatten := map[string]interface{}{
		"id":              s3_storage.ID,
		"name":            s3_storage.Name,
		"backend":         s3_storage.Backend,
		"client_endpoint": s3_storage.ClientEndpoint,
		"access_key":      s3_storage.AccessKey,
		"secret_key":      s3_storage.SecretKey,
	}

	if err := setResourceDataFromMap(d, flatten); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(s3_storage.ID)
	return nil
}
