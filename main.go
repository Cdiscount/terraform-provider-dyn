package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/dyn"
)
// Generate docs for website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dyn.Provider})
}
