package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceJobStorageRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

	teamID, ok := d.Get("team_id").(int)
	if !ok {
		return fmt.Errorf("team_id is not a int")
	}

	jobStorage, err := paperspaceClient.GetJobStorageByRegion(teamID, paperspaceClient.Region)
	if err != nil {
		return err
	}
	if jobStorage.Handle == "" {
		return errors.New("Could not find job storage")
	}

	d.SetId(jobStorage.Handle)
	updateJobStorageSchema(d, jobStorage)

	return nil
}

func updateJobStorageSchema(d *schema.ResourceData, jobStorage JobStorage) {
	d.Set("handle", jobStorage.Handle)
}

func dataSourceJobStorage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJobStorageRead,
		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"handle": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
