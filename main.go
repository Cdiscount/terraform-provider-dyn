package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/dyn"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dyn.Provider})
}
