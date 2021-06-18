package dyn

import (
	"fmt"

	"github.com/Cdiscount/terraform-provider-dyn/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDynDSFMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDSFMonitorCreate,
		Read:   resourceDynDSFMonitorRead,
		Update: resourceDynDSFMonitorUpdate,
		Delete: resourceDynDSFMonitorDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A label to identify the Monitor",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HTTP", "HTTPS", "PING", "SMTP", "TCP"}, false),
				Description: `The protocol to monitor
Valid values:
  * HTTP
  * HTTPS
  * PING
  * SMTP
  * TCP`,
			},
			"response_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 3),
				Description: `Minimum required ‘up’ agent responses to report response pool host as ‘up’. If ‘up’ responses are less than the minimum, host is set to failover.

Valid values: 0, 1, 2 or 3.`,
			},
			"probe_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{60, 300, 600, 900}),
				Description: `How often to run the monitor. Must be twice the TTL setting.
Valid values:
  * 60 – Every minute
  * 300 – Every 5 minutes
  * 600 – Every 10 minutes
  * 900 – Every 15 minutes`,
			},
			"retries": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 2),
				Description: `How many retries the monitor should attempt on failure before giving up.
Valid values:
  * 0
  * 1
  * 2`,
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if the Monitor is active",
			},
			"options": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "Options that pertain to the Monitor",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Time (in seconds) before the connection attempt times out",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "For HTTP(S)/SMTP/TCP probes, an alternate connection port. Leaving the field blank means it will monitor the default port (80 for HTTP and TCP, 443 for HTTPS, and 25 for SMTP)",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "For HTTP(S) probes, a specific path to request. Designate a path other than the root to be monitored. Paths should be supplied as a relative path to the root ‘/’ directory of the website.",
						},
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "For HTTP(S) probes, a value to pass in to the Host: header.",
						},
						"header": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							Description: `For HTTP(S) probes, additional header fields/values to pass in, separated by the newline character (\n).
See [Configuring Monitor Headers](https://help.dyn.com/configuring-monitor-headers/) for more information on using custom headers and macros in your endpoint monitoring.`,
						},
						"expected": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Designate the data expected in the protocol response while monitoring the host in the pool. Maximum length: 255 bytes. Exceeding the maximum size will result in an ‘Invalid_Data’ error at run time with the message ‘Too long’. Field is case-sensitive. Exact string match required to return ‘up’ status. For HTTP(S) probes, a case sensitive sub-string to search for in the response. For SMTP probes, a string to compare the banner against. Not used for PING, or TCP protocols.`,
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
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)

	err = client.Do("POST", "DSFMonitor", request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_monitor(d, &response.Data)

	return nil
}

func resourceDynDSFMonitorRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)
	response := &api.DSFMonitorResponse{}

	url := fmt.Sprintf("DSFMonitor/%s", id)
	err = client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_monitor(d, &response.Data)

	return nil
}

func resourceDynDSFMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)
	request := createRequest(d)
	response := &api.DSFMonitorResponse{}

	url := fmt.Sprintf("DSFMonitor/%s", id)
	err = client.Do("PUT", url, request, response)
	if err != nil {
		return err
	}

	load_dsf_monitor(d, &response.Data)

	return nil
}

func resourceDynDSFMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return err
	}
	defer provider.PutClient(client)

	url := fmt.Sprintf("DSFMonitor/%s", id)
	err = client.Do("DELETE", url, nil, nil)
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
