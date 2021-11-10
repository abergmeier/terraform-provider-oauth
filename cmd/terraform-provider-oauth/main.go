package main

import (
	"log"

	"github.com/abergmeier/terraform-provider-oauth/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {

	// Since terraform catches our output
	// no need for us to have timestamp
	log.SetFlags(0)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})

}
