package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDSFRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDSFRecordSetCreate,
		Read:   resourceDynDSFRecordSetRead,
		Update: resourceDynDSFRecordSetUpdate,
		Delete: resourceDynDSFRecordSetDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_director_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"response_pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dsf_rsfc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rdata_class": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"automation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"serve_count": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fail_count": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"trouble_count": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"eligible": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDynDSFRecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	request := &api.DSFRecordSetRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Label:          d.Get("label").(string),
		RDataClass:     d.Get("rdata_class").(string),
		TTL:            d.Get("ttl").(string),
		Automation:     d.Get("automation").(string),
		ServeCount:     d.Get("serve_count").(string),
		FailCount:      d.Get("fail_count").(string),
		TroubleCount:   d.Get("trouble_count").(string),
		Eligible:       d.Get("eligible").(string),
		MonitorID:      d.Get("monitor_id").(string),
		DSFRsfc:        d.Get("dsf_rsfc_id").(string),
		ResponsePoolId: d.Get("response_pool_id").(string),
	}
	traffic_director_id := d.Get("traffic_director_id").(string)
	response := &api.DSFRecordSetResponse{}
	client := meta.(*api.ConvenientClient)

	url := fmt.Sprintf("DSFRecordSet/%s", traffic_director_id)
	err := client.Do("POST", url, request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_record_set(d, &response.Data)

	return nil
}

func resourceDynDSFRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	response := &api.DSFRecordSetResponse{}

	url := fmt.Sprintf("DSFRecordSet/%s/%s", traffic_director_id, id)
	err := client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_record_set(d, &response.Data)

	return nil
}

func resourceDynDSFRecordSetUpdate(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	request := &api.DSFRecordSetRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Label:        d.Get("label").(string),
		RDataClass:   d.Get("rdata_class").(string),
		TTL:          d.Get("ttl").(string),
		Automation:   d.Get("automation").(string),
		ServeCount:   d.Get("serve_count").(string),
		FailCount:    d.Get("fail_count").(string),
		TroubleCount: d.Get("trouble_count").(string),
		Eligible:     d.Get("eligible").(string),
		MonitorID:    d.Get("monitor_id").(string),
	}
	response := &api.DSFRecordSetResponse{}

	url := fmt.Sprintf("DSFRecordSet/%s/%s", traffic_director_id, id)
	err := client.Do("PUT", url, request, response)
	if err != nil {
		return err
	}

	load_dsf_record_set(d, &response.Data)

	return nil
}

func resourceDynDSFRecordSetDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)

	traffic_director_id := d.Get("traffic_director_id").(string)
	publish := api.PublishBlock{
		Publish: true,
	}
	url := fmt.Sprintf("DSFRecordSet/%s/%s", traffic_director_id, id)
	err := client.Do("DELETE", url, publish, nil)
	if err != nil {
		return err
	}

	return nil
}

func load_dsf_record_set(d *schema.ResourceData, response *api.DSFRecordSet) {
	d.Set("label", response.Label)
	d.Set("automation", response.Automation)
	d.Set("rdata_class", response.RDataClass)
	d.Set("ttl", response.TTL)
	d.Set("serve_count", response.ServeCount)
	d.Set("fail_count", response.FailCount)
	d.Set("trouble_count", response.TroubleCount)
	d.Set("eligible", response.Eligible)
	d.Set("monitor_id", response.MonitorID)
}
