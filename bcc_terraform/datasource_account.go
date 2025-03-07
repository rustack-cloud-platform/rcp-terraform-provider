package bcc_terraform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the current user",
			},
			"email": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of current user",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The username of current user",
			},
		},
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	manager := meta.(*CombinedConfig).Manager()
	account, err := manager.GetAccount()

	if err != nil {
		return diag.Errorf("Error retrieving account: %s", err)
	}

	d.SetId(account.ID)
	d.Set("email", account.Email)
	d.Set("username", account.Username)

	return nil
}
