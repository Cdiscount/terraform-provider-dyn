package dyn

import (
	"fmt"

	"github.com/Cdiscount/terraform-provider-dyn/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDynTrafficDirector() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorCreate,
		Read:   resourceDynTrafficDirectorRead,
		Update: resourceDynTrafficDirectorUpdate,
		Delete: resourceDynTrafficDirectorDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Traffic Director service",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The default TTL to be used across the service",
			},
			"node": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the zone",
						},
						"fqdn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Fully qualified domain name of a node in the zone",
						},
					},
				},
			},
		},
	}
}

func resourceDynTrafficDirectorCreate(d *schema.ResourceData, meta interface{}) error {
	request := &api.DSFServiceRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Label: d.Get("label").(string),
		TTL:   api.SInt(d.Get("ttl").(int)),
	}
	response := &api.DSFResponse{}

	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)

	err = client.Do("POST", "DSF", request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_service(d, &response.Data)

	return updateDsfNodes(d, client)
}

func resourceDynTrafficDirectorRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)

	response := &api.DSFResponse{}

	url := fmt.Sprintf("DSF/%s", id)
	err = client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_service(d, &response.Data)
	load_nodes(response.Data.Nodes, d)

	return nil
}

func resourceDynTrafficDirectorUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)

	if d.HasChanges("label", "ttl") {
		id := d.Id()
		request := &api.DSFServiceRequest{
			PublishBlock: api.PublishBlock{
				Publish: true,
			},
			Label: d.Get("label").(string),
			TTL:   api.SInt(d.Get("ttl").(int)),
		}
		response := &api.DSFResponse{}

		url := fmt.Sprintf("DSF/%s", id)
		err := client.Do("PUT", url, request, response)
		if err != nil {
			return err
		}
		load_dsf_service(d, &response.Data)
	}
	return updateDsfNodes(d, client)
}

func resourceDynTrafficDirectorDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)

	publish := api.PublishBlock{
		Publish: true,
	}
	url := fmt.Sprintf("DSF/%s", id)
	err = client.Do("DELETE", url, &publish, nil)
	if err != nil {
		return err
	}

	return nil
}

func updateDsfNodes(d *schema.ResourceData, client *api.ConvenientClient) error {
	id := d.Id()
	request := &api.DSFNodeRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Node: nodes_from_schema(d),
	}
	response := &api.DSFNodeResponse{}
	url_node := fmt.Sprintf("DSFNode/%s", id)

	err := client.Do("PUT", url_node, request, response)
	if err != nil {
		return err
	}
	load_nodes(response.Data, d)
	return nil
}

func nodes_from_schema(d *schema.ResourceData) []api.DSFNode {
	raw_nodes := d.Get("node").([]interface{})
	nodes := make([]api.DSFNode, len(raw_nodes))
	for i, i_node := range raw_nodes {
		node := i_node.(map[string]interface{})
		nodes[i] = api.DSFNode{
			Zone: node["zone"].(string),
			FQDN: node["fqdn"].(string),
		}
	}
	return nodes
}

func load_dsf_service(d *schema.ResourceData, response *api.DSFService) {
	d.Set("label", response.Label)
	d.Set("ttl", response.TTL)
}

func load_nodes(raw_nodes []api.DSFNode, d *schema.ResourceData) {
	nodes := make([]map[string]interface{}, len(raw_nodes))
	for i, node := range raw_nodes {
		nodes[i] = map[string]interface{}{
			"zone": node.Zone,
			"fqdn": node.FQDN,
		}
	}
	d.Set("node", nodes)
}
