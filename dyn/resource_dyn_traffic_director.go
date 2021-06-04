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
		},
	}
}

func resourceDynTrafficDirectorCreate(d *schema.ResourceData, meta interface{}) error {
	request := &api.DSFServiceRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
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

	return nil
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
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	request := &api.DSFServiceRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
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

	return nil
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

func load_dsf_service(d *schema.ResourceData, response *api.DSFService) {
	d.Set("label", response.Label)
	d.Set("ttl", response.TTL)
}
