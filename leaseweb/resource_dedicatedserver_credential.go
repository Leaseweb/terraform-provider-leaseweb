package leaseweb

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDedicatedServerCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedServerCredentialCreate,
		ReadContext:   resourceDedicatedServerCredentialRead,
		UpdateContext: resourceDedicatedServerCredentialUpdate,
		DeleteContext: resourceDedicatedServerCredentialDelete,
		Schema: map[string]*schema.Schema{
			"dedicated_server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"}, false),
				ForceNew:     true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.SplitN(d.Id(), ":", 3)

				if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
					return nil, fmt.Errorf("Invalid ID format (%s), expected dedicated_server_id:credential_type:credential_username", d.Id())
				}

				d.Set("dedicated_server_id", parts[0])
				d.Set("type", parts[1])
				d.Set("username", parts[2])
				d.SetId(parts[0] + parts[1] + parts[2])

				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDedicatedServerCredentialCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var credential = Credential{
		Type:     d.Get("type").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	createdCredential, err := createDedicatedServerCredential(serverID, &credential)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serverID + createdCredential.Type + createdCredential.Username)

	return resourceDedicatedServerCredentialRead(ctx, d, m)
}

func resourceDedicatedServerCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	credentialType := d.Get("type").(string)
	username := d.Get("username").(string)

	var diags diag.Diagnostics

	credential, err := getDedicatedServerCredential(serverID, credentialType, username)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("type", credential.Type)
	d.Set("username", credential.Username)
	d.Set("password", credential.Password)

	return diags
}

func resourceDedicatedServerCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)

	var credential = Credential{
		Type:     d.Get("type").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	if _, err := updateDedicatedServerCredential(serverID, &credential); err != nil {
		return diag.FromErr(err)
	}

	return resourceDedicatedServerCredentialRead(ctx, d, m)
}

func resourceDedicatedServerCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	serverID := d.Get("dedicated_server_id").(string)

	var credential = Credential{
		Type:     d.Get("type").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	if err := deleteDedicatedServerCredential(serverID, &credential); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
