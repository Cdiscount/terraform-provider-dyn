package dyn

import (
	"sync"

	"github.com/Cdiscount/terraform-provider-dyn/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	schema.DescriptionKind = schema.StringMarkdown
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
			"dyn_dsf_monitor":       resourceDynDSFMonitor(),
		},

		ConfigureFunc: providerConfigure,
	}
}

type DynProvider struct {
	config  *Config
	clients []*api.ConvenientClient
	mutex   sync.Mutex
}

// Get a client from the pool, creating a new one if necessary
func (p *DynProvider) GetClient() (*api.ConvenientClient, error) {
	p.mutex.Lock()
	if len(p.clients) > 0 {
		client := p.clients[len(p.clients)-1]
		p.clients = p.clients[:len(p.clients)-1]
		p.mutex.Unlock()
		return client, nil
	}
	p.mutex.Unlock()
	return p.config.Client()
}

// Put back a client to the pool.
// If not done, the client is lost and a new one must be created
func (p *DynProvider) PutClient(c *api.ConvenientClient) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.clients = append(p.clients, c)
}

func GetProvider(meta interface{}) *DynProvider {
	return meta.(*DynProvider)
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		CustomerName: d.Get("customer_name").(string),
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
	}

	provider := DynProvider{
		config:  &config,
		clients: make([]*api.ConvenientClient, 0, 10),
	}
	return &provider, nil
}
