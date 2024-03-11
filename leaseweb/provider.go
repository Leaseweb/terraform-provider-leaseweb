package leaseweb

import (
	"context"
	"runtime"

	LSW "github.com/LeaseWeb/leaseweb-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// VERSION is a placeholder for the actual version tag which is set at buildtime via compiler flags
var VERSION = "dev"

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Description: `
The base URL of the API endpoint to use.
By default it takes the value from the ` + "`LEASEWEB_API_URL`" + ` environment variable if present,
otherwise it defaults to "https://api.leaseweb.com".
`,
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEASEWEB_API_URL", nil),
			},
			"api_token": {
				Description: `
The API token to use.
By default it takes the value from the ` + "`LEASEWEB_API_TOKEN`" + ` environment variable if present.
`,
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEASEWEB_API_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"leaseweb_dedicated_server_credential":                       resourceDedicatedServerCredential(),
			"leaseweb_dedicated_server_installation":                     resourceDedicatedServerInstallation(),
			"leaseweb_dedicated_server_notification_setting_bandwidth":   resourceDedicatedServerNotificationSettingBandwidth(),
			"leaseweb_dedicated_server_notification_setting_datatraffic": resourceDedicatedServerNotificationSettingDatatraffic(),
			"leaseweb_dedicated_server":                                  resourceDedicatedServer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"leaseweb_dedicated_server_control_panels":    dataSourceDedicatedServerControlPanels(),
			"leaseweb_dedicated_server_credential":        dataSourceDedicatedServerCredential(),
			"leaseweb_dedicated_server_operating_systems": dataSourceDedicatedServerOperatingSystems(),
			"leaseweb_dedicated_servers":                  dataSourceDedicatedServers(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	baseURL := d.Get("api_url").(string)
	apiToken := d.Get("api_token").(string)

	if apiToken == "" {
		return nil, diag.Errorf("missing leaseweb provider token")
	}

	var diags diag.Diagnostics

	LSW.InitLeasewebClient(apiToken)
	if baseURL != "" {
		LSW.SetBaseUrl(baseURL)
	}

	LSW.SetUserAgent("terraform-provider-leaseweb/" + VERSION + " (" + runtime.GOOS + "; " + runtime.GOARCH + "; " + runtime.Version() + ")")

	return nil, diags
}
