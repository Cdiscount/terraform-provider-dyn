package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDSFResponsePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDSFResponsePoolCreate,
		Read:   resourceDynDSFResponsePoolRead,
		Update: resourceDynDSFResponsePoolUpdate,
		Delete: resourceDynDSFResponsePoolDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"automation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"traffic_director_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDynDSFResponsePoolCreate(d *schema.ResourceData, meta interface{}) error {
	request := &api.DSFResponsePoolRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
			Publish: true,
		},
		Label:      d.Get("label").(string),
		Automation: d.Get("automation").(string),
	}
	traffic_director_id := d.Get("traffic_director_id").(string)
	response := &api.DSFResponsePoolResponse{}
	client := meta.(*api.ConvenientClient)

	url := fmt.Sprintf("DSFResponsePool/%s", traffic_director_id)
	err := client.Do("POST", url, request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_response_pool(d, &response.Data)

	return nil
}

func resourceDynDSFResponsePoolRead(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	response := &api.DSFResponsePoolResponse{}

	url := fmt.Sprintf("DSFResponsePool/%s/%s", traffic_director_id, id)
	err := client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_response_pool(d, &response.Data)

	return nil
}

func resourceDynDSFResponsePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	request := &api.DSFResponsePoolRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
			Publish: true,
		},
		Label:      d.Get("label").(string),
		Automation: d.Get("automation").(string),
	}
	response := &api.DSFResponsePoolResponse{}

	url := fmt.Sprintf("DSFResponsePool/%s/%s", traffic_director_id, id)
	err := client.Do("PUT", url, request, response)
	if err != nil {
		return err
	}

	load_dsf_response_pool(d, &response.Data)

	return nil
}

func resourceDynDSFResponsePoolDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)

	traffic_director_id := d.Get("traffic_director_id").(string)
	publish := api.CreateOrUpdateBlock{
		Publish: true,
	}
	url := fmt.Sprintf("DSFResponsePool/%s/%s", traffic_director_id, id)
	err := client.Do("DELETE", url, publish, nil)
	if err != nil {
		return err
	}

	return nil
}

func load_dsf_response_pool(d *schema.ResourceData, response *api.DSFResponsePool) {
	d.Set("label", response.Label)
	d.Set("automation", response.Automation)
}
