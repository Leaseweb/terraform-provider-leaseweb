package leaseweb

import (
	"context"
	"strconv"
	"time"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDedicatedServerCredential() *schema.Resource {
	return &schema.Resource{
		Description: `
The ` + "`dedicated_server_credential`" + ` data source allows access to the list of
credentials available for a dedicated server.
`,
		ReadContext: dataSourceDedicatedServerCredentialRead,
		Schema: map[string]*schema.Schema{
			"dedicated_server_id": {
				Description: "The ID of the dedicated server.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"type": {
				Description: `
The type of the credential.
Can be either ` + "`OPERATING_SYSTEM`" + `, ` + "`CONTROL_PANEL`" + `, ` + "`REMOTE_MANAGEMENT`" + `, ` + "`RESCUE_MODE`" + `, ` + "`SWITCH`" + `, ` + "`PDU`" + `, ` + "`FIREWALL`" + ` or ` + "`LOAD_BALANCER`" + `.
`,
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"}, false),
			},
			"username": {
				Description: "The username of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The password of the credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDedicatedServerCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	serverID := d.Get("dedicated_server_id").(string)
	credentialType := d.Get("type").(string)
	username := d.Get("username").(string)

	credential, err := LSW.DedicatedServerApi{}.GetCredential(ctx, serverID, credentialType, username)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("password", credential.Password)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
