package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDSFRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDSFRecordSetCreate,
		Read:   resourceDynDSFRecordSetRead,
		Update: resourceDynDSFRecordSetUpdate,
		Delete: resourceDynDSFRecordSetDelete,

		Description: "Dynect traffic director record set",
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Record set name",
			},
			"traffic_director_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The traffic director in which the record set is created",
			},
			"response_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The response pool in which the record set is created",
			},
			"dsf_rsfc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identifier for the Response Pool whose Record Set Failover Chain will include this Record Set",
			},
			"rdata_class": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA"}, false),
				Description:  "The type of rdata represented by this Record Set",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Default TTL used for Records within this Record Set",
			},
			"automation": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"auto", "auto_down", "manual"}, false),
				Description: `Defines how eligible can be changed in response to monitoring
  * auto — Sets the serve_mode field to ‘Monitor & Obey’. Default
  * auto_down — Sets the serve_mode field to ‘Monitor & Remove’
  * manual — Couples with eligible value to determine other serve_mode field values`,
			},
			"serve_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "How many Records to serve out of this Record Set",
			},
			"fail_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of Records that must not be okay before the Record Set becomes ineligible",
			},
			"trouble_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of Records that must not be okay before the Record Set becomes in trouble",
			},
			"eligible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: `Indicates whether or not the Record Set can be served
  * false — When automation is set to manual, sets the serve_mode field to ‘Do Not Serve’
  * true — Default. When automation is set to manual, Record Set can be served`,
			},
			"monitor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The id of the monitoring object",
			},
		},
	}
}

func resourceDynDSFRecordSetCreate(d *schema.ResourceData, meta interface{}) error {
	request := computeDSFRecordSetRequest(d, true)
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
	request := computeDSFRecordSetRequest(d, false)
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

func computeDSFRecordSetRequest(d *schema.ResourceData, isCreate bool) *api.DSFRecordSetRequest {
	request := &api.DSFRecordSetRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Label:        d.Get("label").(string),
		RDataClass:   d.Get("rdata_class").(string),
		TTL:          api.SInt(d.Get("ttl").(int)),
		Automation:   d.Get("automation").(string),
		ServeCount:   api.SInt(d.Get("serve_count").(int)),
		FailCount:    api.SInt(d.Get("fail_count").(int)),
		TroubleCount: api.SInt(d.Get("trouble_count").(int)),
		Eligible:     api.SBool(d.Get("eligible").(bool)),
		MonitorID:    d.Get("monitor_id").(string),
		DSFRsfc:      d.Get("dsf_rsfc_id").(string),
	}
	if isCreate {
		request.ResponsePoolId = d.Get("response_pool_id").(string)
	}
	return request
}

func load_dsf_record_set(d *schema.ResourceData, response *api.DSFRecordSet) {
	d.Set("label", response.Label)
	d.Set("automation", response.Automation)
	d.Set("rdata_class", response.RDataClass)
	d.Set("ttl", response.TTL)
	d.Set("serve_count", response.ServeCount)
	d.Set("fail_count", response.FailCount)
	d.Set("trouble_count", response.TroubleCount)
	d.Set("eligible", bool(response.Eligible))
	d.Set("monitor_id", response.MonitorID)
}
