package leaseweb

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"leaseweb_api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEASEWEB_API_URL", "https://api.leaseweb.com"),
			},
			"leaseweb_api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LEASEWEB_API_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"leaseweb_dedicatedserver": resourceDedicatedServer(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	baseURL := d.Get("leaseweb_api_url").(string)
	apiToken := d.Get("leaseweb_api_token").(string)

	if baseURL == "" || apiToken == "" {
		return nil, diag.Errorf("missing leaseweb provider base url or token")
	}

	var diags diag.Diagnostics

	leasewebAPIURL = baseURL
	leasewebAPIToken = apiToken
	leasewebClient = &http.Client{Timeout: 10 * time.Second}

	return nil, diags
}
