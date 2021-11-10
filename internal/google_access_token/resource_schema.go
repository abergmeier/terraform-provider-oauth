package google_access_token

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"scopes": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Default: []string{"https://www.googleapis.com/auth/cloud-platform"},
			},
			"access_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"token_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		ReadContext: read,
	}
}
