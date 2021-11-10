package google_access_token

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
