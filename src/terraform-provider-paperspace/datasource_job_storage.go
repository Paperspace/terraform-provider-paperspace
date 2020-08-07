package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceJobStorageRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := newInternalPaperspaceClient(m)
	region := paperspaceClient.Region

	teamID, ok := d.Get("team_id").(int)
	if !ok {
		return fmt.Errorf("team_id is not a int")
	}
	regionData, ok := d.Get("region").(string)
	if !ok {
		return fmt.Errorf("region is not a string")
	}

	if regionData != "" {
		region = regionData
	}

	jobStorage, err := paperspaceClient.GetJobStorageByRegion(teamID, region)
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
