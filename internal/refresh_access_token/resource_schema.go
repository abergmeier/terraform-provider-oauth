package refresh_access_token

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client ID",
				Sensitive:   true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Client Secret",
				Sensitive:   true,
			},
			"refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Refresh Token",
				Sensitive:   true,
			},
			"token_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://oauth2.googleapis.com/token",
			},
			"access_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"id_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"token_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: read,
	}
}
