package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDsfRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDsfRecordCreate,
		Read:   resourceDynDsfRecordRead,
		Update: resourceDynDsfRecordUpdate,
		Delete: resourceDynDsfRecordDelete,

		Schema: map[string]*schema.Schema{
			"record_set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_director_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"automation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"master_line": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDynDsfRecordCreate(d *schema.ResourceData, meta interface{}) error {
	request := &api.DSFRecordRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
			Publish: true,
		},
		Label:      d.Get("label").(string),
		Weight:     d.Get("weight").(string),
		Automation: d.Get("automation").(string),
		MasterLine: d.Get("master_line").(string),
	}

	traffic_director_id := d.Get("traffic_director_id").(string)
	record_set_id := d.Get("record_set_id").(string)
	url := fmt.Sprintf("DSFRecord/%s/%s", traffic_director_id, record_set_id)

	response := &api.DSFRecordResponse{}
	client := meta.(*api.ConvenientClient)
	err := client.Do("POST", url, request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_record(d, &response.Data)

	return nil
}

func resourceDynDsfRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.ConvenientClient)
	response := &api.DSFRecordResponse{}

	id := d.Id()
	traffic_director_id := d.Get("traffic_director_id").(string)
	url := fmt.Sprintf("DSFRecord/%s/%s", traffic_director_id, id)

	err := client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_record(d, &response.Data)

	return nil
}

func resourceDynDsfRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.ConvenientClient)
	request := &api.DSFRecordRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
			Publish: true,
		},
		Label:      d.Get("label").(string),
		Weight:     d.Get("weight").(string),
		Automation: d.Get("automation").(string),
		MasterLine: d.Get("master_line").(string),
	}
	response := &api.DSFRecordResponse{}

	id := d.Id()
	traffic_director_id := d.Get("traffic_director_id").(string)
	url := fmt.Sprintf("DSFRecord/%s/%s", traffic_director_id, id)

	err := client.Do("PUT", url, request, response)
	if err != nil {
		return err
	}

	load_dsf_record(d, &response.Data)

	return nil
}

func resourceDynDsfRecordDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)

	url := fmt.Sprintf("DSFRecord/%s", id)
	err := client.Do("DELETE", url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func load_dsf_record(d *schema.ResourceData, response *api.DSFRecord) {
	d.Set("label", response.Label)
	d.Set("weight", response.Weight)
	d.Set("automation", response.Automation)
	d.Set("master_line", response.MasterLine)
}
