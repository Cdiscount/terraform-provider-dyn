package main

import (
	"github.com/Cdiscount/terraform-provider-dyn/dyn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Generate docs for website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dyn.Provider})
}
