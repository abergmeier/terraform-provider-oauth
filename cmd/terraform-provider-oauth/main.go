package main

import (
	"github.com/abergmeier/terraform-provider-oauth/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})

}
