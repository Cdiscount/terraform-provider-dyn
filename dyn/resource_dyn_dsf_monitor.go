package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDSFMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDSFMonitorCreate,
		Read:   resourceDynDSFMonitorRead,
		Update: resourceDynDSFMonitorUpdate,
		Delete: resourceDynDSFMonitorDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "PING", "SMTP", "TCP"}, false),
			},
			"response_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 3),
			},
			"probe_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{60, 300, 600, 900}),
			},
			"retries": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 2),
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"header": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"expected": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDynDSFMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	request := createRequest(d)
	response := &api.DSFMonitorResponse{}
	client := meta.(*api.ConvenientClient)

	err := client.Do("POST", "DSFMonitor", request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_monitor(d, &response.Data)

	return nil
}

func resourceDynDSFMonitorRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	response := &api.DSFMonitorResponse{}

	url := fmt.Sprintf("DSFMonitor/%s", id)
	err := client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_monitor(d, &response.Data)

	return nil
}

func resourceDynDSFMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	request := createRequest(d)
	response := &api.DSFMonitorResponse{}

	url := fmt.Sprintf("DSFMonitor/%s", id)
	err := client.Do("PUT", url, request, response)
	if err != nil {
		return err
	}

	load_dsf_monitor(d, &response.Data)

	return nil
}

func resourceDynDSFMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)

	url := fmt.Sprintf("DSFMonitor/%s", id)
	err := client.Do("DELETE", url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func createRequest(d *schema.ResourceData) *api.DSFMonitor {
	var options *api.DSFMonitorOptions
	var raw_options = d.Get("options").([]interface{})
	if len(raw_options) > 0 {
		opts := raw_options[0].(map[string]interface{})
		options = &api.DSFMonitorOptions{
			Timeout:  api.SInt(opts["timeout"].(int)),
			Port:     api.SInt(opts["port"].(int)),
			Path:     opts["path"].(string),
			Host:     opts["host"].(string),
			Header:   opts["header"].(string),
			Expected: opts["expected"].(string),
		}
	}
	request := &api.DSFMonitor{
		Label:         d.Get("label").(string),
		Protocol:      d.Get("protocol").(string),
		ResponseCount: api.SInt(d.Get("response_count").(int)),
		ProbeInterval: api.SInt(d.Get("probe_interval").(int)),
		Retries:       api.SInt(d.Get("retries").(int)),
		Active:        api.YNBool(d.Get("active").(bool)),
		Options:       options,
	}
	return request
}

func load_dsf_monitor(d *schema.ResourceData, response *api.DSFMonitor) {
	d.Set("label", response.Label)
	d.Set("protocol", response.Protocol)
	d.Set("response_count", response.ResponseCount)
	d.Set("probe_interval", response.ProbeInterval)
	d.Set("retries", response.Retries)
	d.Set("active", response.Active)
	options := make([]map[string]interface{}, 0, 1)
	if response.Options != nil {
		option := map[string]interface{}{
			"timeout":  response.Options.Timeout,
			"port":     response.Options.Port,
			"path":     response.Options.Path,
			"host":     response.Options.Host,
			"header":   response.Options.Header,
			"expected": response.Options.Expected,
		}
		options = append(options, option)
	}
	d.Set("options", options)
}
