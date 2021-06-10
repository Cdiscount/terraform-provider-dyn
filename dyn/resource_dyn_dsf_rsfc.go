package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDSFRsfc() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDSFRsfcCreate,
		Read:   resourceDynDSFRsfcRead,
		Update: resourceDynDSFRsfcUpdate,
		Delete: resourceDynDSFRsfcDelete,

		Description: "Dynect RecordSet Failover Chain",
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A label for the Record Set Failover Chain",
			},
			"traffic_director_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The traffic director in which we create the ressource",
			},
			"response_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The response pool id in which we create the ressource",
			},
		},
	}
}

func resourceDynDSFRsfcCreate(d *schema.ResourceData, meta interface{}) error {
	request := &api.DSFRsfcRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Label: d.Get("label").(string),
	}
	traffic_director_id := d.Get("traffic_director_id").(string)
	response_pool_id := d.Get("response_pool_id").(string)
	response := &api.DSFRsfcResponse{}
	client := meta.(*api.ConvenientClient)

	url := fmt.Sprintf("DSFRecordSetFailoverChain/%s/%s", traffic_director_id, response_pool_id)
	err := client.Do("POST", url, request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_rsfc(d, &response.Data)

	return nil
}

func resourceDynDSFRsfcRead(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	response := &api.DSFRsfcResponse{}

	url := fmt.Sprintf("DSFRecordSetFailoverChain/%s/%s", traffic_director_id, id)
	err := client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_rsfc(d, &response.Data)

	return nil
}

func resourceDynDSFRsfcUpdate(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	request := &api.DSFRsfcRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Label: d.Get("label").(string),
	}
	response := &api.DSFRsfcResponse{}

	url := fmt.Sprintf("DSFRecordSetFailoverChain/%s/%s", traffic_director_id, id)
	err := client.Do("PUT", url, request, response)
	if err != nil {
		return err
	}

	load_dsf_rsfc(d, &response.Data)

	return nil
}

func resourceDynDSFRsfcDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)

	traffic_director_id := d.Get("traffic_director_id").(string)
	publish := api.PublishBlock{
		Publish: true,
	}
	url := fmt.Sprintf("DSFRecordSetFailoverChain/%s/%s", traffic_director_id, id)
	err := client.Do("DELETE", url, publish, nil)
	if err != nil {
		return err
	}

	return nil
}

func load_dsf_rsfc(d *schema.ResourceData, response *api.DSFRecordSetChain) {
	d.Set("label", response.Label)
}
