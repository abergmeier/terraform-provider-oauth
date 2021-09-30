package provider

import (
	"github.com/abergmeier/terraform-provider-oauth/internal/refresh_access_token"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"oauth_refresh_access_token": refresh_access_token.Resource(),
		},
		ResourcesMap: map[string]*schema.Resource{},
		Schema:       map[string]*schema.Schema{},
	}
}
