package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-dyn/dyn"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dyn.Provider})
}
