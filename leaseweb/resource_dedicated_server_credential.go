package leaseweb

import (
	"context"
	"fmt"
	"strings"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDedicatedServerCredential() *schema.Resource {
	return &schema.Resource{
		Description: `
The ` + "`dedicated_server_credential`" + ` resource manages a credential
linked to a dedicated server.
`,
		CreateContext: resourceDedicatedServerCredentialCreate,
		ReadContext:   resourceDedicatedServerCredentialRead,
		UpdateContext: resourceDedicatedServerCredentialUpdate,
		DeleteContext: resourceDedicatedServerCredentialDelete,
		Schema: map[string]*schema.Schema{
			"dedicated_server_id": {
				Description: "The ID of the dedicated server.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: `
The type of the credential.
Can be either ` + "`OPERATING_SYSTEM`" + `, ` + "`CONTROL_PANEL`" + `, ` + "`REMOTE_MANAGEMENT`" + `, ` + "`RESCUE_MODE`" + `, ` + "`SWITCH`" + `, ` + "`PDU`" + `, ` + "`FIREWALL`" + ` or ` + "`LOAD_BALANCER`" + `.
`,
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"}, false),
				ForceNew:     true,
			},
			"username": {
				Description: "The username of the credential.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Description: "The password of the credential.",
				Type:        schema.TypeString,
				Required:    true,
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
	credentialType := d.Get("type").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	createdCredential, err := LSW.DedicatedServerApi{}.CreateCredential(ctx, serverID, credentialType, username, password)
	if err != nil {
		logSdkApiError(ctx, err)
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

	credential, err := LSW.DedicatedServerApi{}.GetCredential(ctx, serverID, credentialType, username)
	if err != nil {
		logSdkApiError(ctx, err)
		return diag.FromErr(err)
	}

	d.Set("password", credential.Password)

	return diags
}

func resourceDedicatedServerCredentialUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	serverID := d.Get("dedicated_server_id").(string)
	credentialType := d.Get("type").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	_, err := LSW.DedicatedServerApi{}.UpdateCredential(ctx, serverID, credentialType, username, password)
	if err != nil {
		logSdkApiError(ctx, err)
		return diag.FromErr(err)
	}

	return resourceDedicatedServerCredentialRead(ctx, d, m)
}

func resourceDedicatedServerCredentialDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	serverID := d.Get("dedicated_server_id").(string)
	credentialType := d.Get("type").(string)
	username := d.Get("username").(string)

	err := LSW.DedicatedServerApi{}.DeleteCredential(ctx, serverID, credentialType, username)
	if err != nil {
		logSdkApiError(ctx, err)
		return diag.FromErr(err)
	}

	return diags
}
