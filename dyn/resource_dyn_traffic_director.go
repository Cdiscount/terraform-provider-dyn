package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynTrafficDirector() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynTrafficDirectorCreate,
		Read:   resourceDynTrafficDirectorRead,
		Update: resourceDynTrafficDirectorUpdate,
		Delete: resourceDynTrafficDirectorDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"node": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fqdn": {
							Type:     schema.TypeString,
							Required: true,
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
		TTL:   d.Get("ttl").(string),
	}
	response := &api.DSFResponse{}
	client := meta.(*api.ConvenientClient)

	err := client.Do("POST", "DSF", request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_service(d, &response.Data)

	return updateDsfNodes(d, client)
}

func resourceDynTrafficDirectorRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	response := &api.DSFResponse{}

	url := fmt.Sprintf("DSF/%s", id)
	err := client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_service(d, &response.Data)

	return nil
}

func resourceDynTrafficDirectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.ConvenientClient)
	if d.HasChanges("label", "ttl") {
		id := d.Id()
		request := &api.DSFServiceRequest{
			PublishBlock: api.PublishBlock{
				Publish: true,
			},
			Label: d.Get("label").(string),
			TTL:   d.Get("ttl").(string),
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
	client := meta.(*api.ConvenientClient)

	url := fmt.Sprintf("DSF/%s", id)
	err := client.Do("DELETE", url, nil, nil)
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
	load_nodes(response.Nodes, d)
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
