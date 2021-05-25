package dyn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynDSFRuleset() *schema.Resource {
	return &schema.Resource{
		Create: resourceDynDSFRulesetCreate,
		Read:   resourceDynDSFRulesetRead,
		Update: resourceDynDSFRulesetUpdate,
		Delete: resourceDynDSFRulesetDelete,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:     schema.TypeString,
				Required: true,
			},
			"traffic_director_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDynDSFRulesetCreate(d *schema.ResourceData, meta interface{}) error {
	request := &api.DSFRulesetRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
			Publish: true,
		},
		Label:        d.Get("label").(string),
		CriteriaType: "always",
	}
	traffic_director_id := d.Get("traffic_director_id").(string)
	response := &api.DSFRulesetResponse{}
	client := meta.(*api.ConvenientClient)

	url := fmt.Sprintf("DSFRuleset/%s", traffic_director_id)
	err := client.Do("POST", url, request, response)
	if err != nil {
		return err
	}

	d.SetId(response.Data.ID)
	load_dsf_ruleset(d, &response.Data)

	return nil
}

func resourceDynDSFRulesetRead(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	response := &api.DSFRulesetResponse{}

	url := fmt.Sprintf("DSFRuleset/%s/%s", traffic_director_id, id)
	err := client.Do("GET", url, nil, response)
	if err != nil {
		return err
	}

	load_dsf_ruleset(d, &response.Data)

	return nil
}

func resourceDynDSFRulesetUpdate(d *schema.ResourceData, meta interface{}) error {
	traffic_director_id := d.Get("traffic_director_id").(string)
	id := d.Id()
	client := meta.(*api.ConvenientClient)
	request := &api.DSFRulesetRequest{
		CreateOrUpdateBlock: api.CreateOrUpdateBlock{
			Publish: true,
		},
		Label:        d.Get("label").(string),
		CriteriaType: "always",
	}
	response := &api.DSFRulesetResponse{}

	url := fmt.Sprintf("DSFRuleset/%s/%s", traffic_director_id, id)
	err := client.Do("PUT", url, request, response)
	if err != nil {
		return err
	}

	load_dsf_ruleset(d, &response.Data)

	return nil
}

func resourceDynDSFRulesetDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	client := meta.(*api.ConvenientClient)

	traffic_director_id := d.Get("traffic_director_id").(string)
	publish := api.CreateOrUpdateBlock{
		Publish: true,
	}
	url := fmt.Sprintf("DSFRuleset/%s/%s", traffic_director_id, id)
	err := client.Do("DELETE", url, publish, nil)
	if err != nil {
		return err
	}

	return nil
}

func load_dsf_ruleset(d *schema.ResourceData, response *api.DSFRuleset) {
	d.Set("label", response.Label)
}
