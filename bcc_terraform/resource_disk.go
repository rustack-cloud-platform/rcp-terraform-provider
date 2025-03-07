package bcc_terraform

import (
	"context"
	"log"
	"time"

	"github.com/basis-cloud/bcc-go/bcc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDisk() *schema.Resource {
	args := Defaults()
	args.injectCreateDisk()
	args.injectContextVdcById()
	args.injectContextStorageProfileById() // override storage_profile_id

	return &schema.Resource{
		CreateContext: resourceDiskCreate,
		ReadContext:   resourceDiskRead,
		UpdateContext: resourceDiskUpdate,
		DeleteContext: resourceDiskDelete,
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

func resourceDiskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	targetVdc, err := GetVdcById(d, manager)
	if err != nil {
		return diag.Errorf("vdc_id: Error getting VDC: %s", err)
	}

	targetStorageProfile, err := GetStorageProfileById(d.Get("storage_profile_id").(string), manager, targetVdc)
	if err != nil {
		return diag.Errorf("storage_profile: Error getting storage profile: %s", err)
	}

	newDisk := bcc.NewDisk(d.Get("name").(string), d.Get("size").(int), targetStorageProfile)
	targetVdc.WaitLock()
	newDisk.Tags = unmarshalTagNames(d.Get("tags"))
	err = targetVdc.CreateDisk(&newDisk)
	if err != nil {
		return diag.Errorf("Error creating disk: %s", err)
	}
	newDisk.WaitLock()

	d.SetId(newDisk.ID)
	log.Printf("[INFO] Disk created, ID: %s", d.Id())

	return resourceDiskRead(ctx, d, meta)
}

func resourceDiskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	disk, err := manager.GetDisk(d.Id())
	if err != nil {
		if err.(*bcc.ApiError).Code() == 404 {
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("id: Error getting disk: %s", err)
		}
	}

	d.SetId(disk.ID)
	d.Set("name", disk.Name)
	d.Set("size", disk.Size)
	d.Set("external_id", disk.ExternalID)
	d.Set("tags", marshalTagNames(disk.Tags))

	return nil
}

func resourceDiskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	disk, err := manager.GetDisk(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting disk: %s", err)
	}

	shouldUpdate := false
	if d.HasChange("name") {
		disk.Name = d.Get("name").(string)
		shouldUpdate = true
	}

	if d.HasChange("tags") {
		disk.Tags = unmarshalTagNames(d.Get("tags"))
		shouldUpdate = true
	}

	if d.HasChange("size") {
		disk.Size = d.Get("size").(int)
		if disk.Locked {
			disk.WaitLock()
		}
		err = disk.Resize(d.Get("size").(int))
		if err != nil {
			return diag.Errorf("size: Error resizing disk: %s", err)
		}
		shouldUpdate = false
	}

	if d.HasChange("storage_profile_id") {
		targetVdc, err := GetVdcById(d, manager)
		if err != nil {
			return diag.Errorf("Error getting VDC: %s", err)
		}

		targetStorageProfileId := d.Get("storage_profile_id").(string)
		targetStorageProfile, err := GetStorageProfileById(targetStorageProfileId, manager, targetVdc)
		if err != nil {
			return diag.Errorf("storage_profile: Error getting storage profile: %s", err)
		}
		if disk.Locked {
			disk.WaitLock()
		}
		err = disk.UpdateStorageProfile(*targetStorageProfile)
		if err != nil {
			return diag.Errorf("storage_profile: Error updating storage: %s", err)
		}
		shouldUpdate = false
	}
	if shouldUpdate {
		if disk.Locked {
			disk.WaitLock()
		}
		disk.Update()
	}

	return resourceDiskRead(ctx, d, meta)
}

func resourceDiskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	disk, err := manager.GetDisk(d.Id())
	if err != nil {
		return diag.Errorf("id: Error getting disk: %s", err)
	}

	if disk.Vm != nil {
		vm, err := manager.GetVm(disk.Vm.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = vm.DetachDisk(disk)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	err = disk.Delete()
	if err != nil {
		return diag.Errorf("Error deleting disk: %s", err)
	}
	disk.WaitLock()

	return nil
}
