package dyn

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDsfRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDynDsfRecordCreate,
		ReadContext:   resourceDynDsfRecordRead,
		UpdateContext: resourceDynDsfRecordUpdate,
		DeleteContext: resourceDynDsfRecordDelete,

		Schema: map[string]*schema.Schema{
			"record_set_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The record set in which this record is added",
			},
			"traffic_director_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The traffic director ID",
			},
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A label for the Record",
			},
			"weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 255),
				Description: `Weight for the Record. Defaults to 1.
  * Valid values for A or AAAA records: 1 – 15.
  * Valid values for CNAME records: 1 – 255.`,
			},
			"automation": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      true,
				ValidateFunc: validation.StringInSlice([]string{"auto", "auto_down", "manual"}, false),
				Description: `Defines how eligible can be changed in response to monitoring.
  * auto — Sets the serve_mode field to ‘Monitor & Obey’. Default.
  * auto_down — Sets the serve_mode field to ‘Monitor & Remove’.
  * manual — Couples with eligible value to determine other serve_mode field values.`,
			},
			"eligible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: `Indicates whether or not the Record can be served.
  * false — When automation is set to manual, sets the serve_mode field to ‘Do Not Serve’.
  * true — Default. When automation is set to manual, sets the serve_mode field to ‘Always Serve’.`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					automation := d.Get("automation").(string)
					return automation != "manual" || old == new
				},
			},
			"master_line": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value to put in the record, i.e. 1.2.3.4 for a DNS A record",
			},
		},
	}
}

func resourceDynDsfRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	request := computeRequest(d)

	traffic_director_id := d.Get("traffic_director_id").(string)
	record_set_id := d.Get("record_set_id").(string)
	url := fmt.Sprintf("DSFRecord/%s/%s", traffic_director_id, record_set_id)

	response := &api.DSFRecordResponse{}
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer provider.PutClient(client)

	err = client.Do("POST", url, request, response)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Data.ID)
	load_dsf_record(d, &response.Data)

	return nil
}

func resourceDynDsfRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer provider.PutClient(client)

	response := &api.DSFRecordResponse{}

	id := d.Id()
	traffic_director_id := d.Get("traffic_director_id").(string)
	url := fmt.Sprintf("DSFRecord/%s/%s", traffic_director_id, id)

	err = client.Do("GET", url, nil, response)
	if err != nil {
		return diag.FromErr(err)
	}

	load_dsf_record(d, &response.Data)

	return nil
}

func resourceDynDsfRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer provider.PutClient(client)

	request := computeRequest(d)
	response := &api.DSFRecordResponse{}

	id := d.Id()
	traffic_director_id := d.Get("traffic_director_id").(string)
	url := fmt.Sprintf("DSFRecord/%s/%s", traffic_director_id, id)

	err = client.Do("PUT", url, request, response)
	if err != nil {
		return diag.FromErr(err)
	}

	load_dsf_record(d, &response.Data)

	return nil
}

func resourceDynDsfRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return diag.FromErr(err)
	}
	defer provider.PutClient(client)

	request := api.PublishBlock{
		Publish: true,
	}
	url := fmt.Sprintf("DSFRecord/%s/%s", traffic_director_id, id)
	err = client.Do("DELETE", url, &request, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func computeRequest(d *schema.ResourceData) *api.DSFRecordRequest {
	request := &api.DSFRecordRequest{
		PublishBlock: api.PublishBlock{
			Publish: true,
		},
		Label:      d.Get("label").(string),
		Weight:     d.Get("weight").(int),
		Automation: d.Get("automation").(string),
		MasterLine: d.Get("master_line").(string),
		Eligible:   nil,
	}
	if request.Automation == "manual" {
		eligible := api.SBool(d.Get("eligible").(bool))
		request.Eligible = &eligible
	}
	return request
}

func load_dsf_record(d *schema.ResourceData, response *api.DSFRecord) {
	d.Set("label", response.Label)
	d.Set("weight", response.Weight)
	d.Set("automation", response.Automation)
	d.Set("master_line", response.MasterLine)
	d.Set("eligible", response.Eligible)
}
