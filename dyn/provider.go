package dyn

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"customer_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYN_CUSTOMER_NAME", nil),
				Description: "A Dyn customer name.",
			},

			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYN_USERNAME", nil),
				Description: "A Dyn username.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DYN_PASSWORD", nil),
				Description: "The Dyn password.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dyn_record":            resourceDynRecord(),
			"dyn_traffic_director":  resourceDynTrafficDirector(),
			"dyn_dsf_ruleset":       resourceDynDSFRuleset(),
			"dyn_dsf_response_pool": resourceDynDSFResponsePool(),
			"dyn_dsf_rsfc":          resourceDynDSFRsfc(),
			"dyn_dsf_record_set":    resourceDynDSFRecordSet(),
			"dyn_dsf_record":        resourceDynDsfRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		CustomerName: d.Get("customer_name").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
	}

	return config.Client()
}
