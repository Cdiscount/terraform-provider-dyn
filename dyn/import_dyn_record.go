package dyn

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"gitlab.cshield.io/cshield.tech/infra/terraform-provider-dyn/api"
)

func resourceDynRecordImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	results := make([]*schema.ResourceData, 1, 1)

	provider := GetProvider(meta)
	client, err := provider.GetClient()
	if err != nil {
		return nil, err
	}
	defer provider.PutClient(client)

	values := strings.Split(d.Id(), "/")

	if len(values) != 3 && len(values) != 4 {
		return nil, fmt.Errorf("invalid id provided, expected format: {type}/{zone}/{fqdn}[/{id}]")
	}

	recordType := values[0]
	recordZone := values[1]
	recordFQDN := values[2]

	var recordID string
	if len(values) == 4 {
		recordID = values[3]
	} else {
		recordID = ""
	}

	record := &api.Record{
		ID:    recordID,
		Name:  "",
		Zone:  recordZone,
		Value: "",
		Type:  recordType,
		FQDN:  recordFQDN,
		TTL:   "",
	}

	// If we already have the record ID, use it for the lookup
	if record.ID == "" {
		err := client.GetRecordID(record)
		if err != nil {
			return nil, err
		}
	} else {
		err := client.GetRecord(record)
		if err != nil {
			return nil, err
		}
	}

	d.SetId(record.ID)
	d.Set("name", record.Name)
	d.Set("zone", record.Zone)
	d.Set("value", record.Value)
	d.Set("type", record.Type)
	d.Set("fqdn", record.FQDN)
	d.Set("ttl", record.TTL)
	results[0] = d

	return results, nil
}
