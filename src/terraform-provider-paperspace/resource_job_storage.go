package main

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceJobStorageCreate(d *schema.ResourceData, m interface{}) error {

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if err := resourceJobStorageRead(d, m); err != nil {
			return resource.RetryableError(err)
		}

		handle, ok := d.Get("handle").(string)
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("handle is not a string"))
		}

		if handle == "" {
			return resource.NonRetryableError(fmt.Errorf("Could not find job storage"))
		}

		return resource.NonRetryableError(nil)
	})
}

func resourceJobStorageRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)
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

	updateJobStorageSchema(d, jobStorage)
	d.SetId(jobStorage.Handle)

	return nil
}

func resourceJobStorageUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO: implement; api doesn't exist yet
	return resourceJobStorageRead(d, m)
}

func resourceJobStorageDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func dataSourceJobStorageRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)
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

func resourceJobStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobStorageCreate,
		Read:   resourceJobStorageRead,
		Update: resourceJobStorageUpdate,
		Delete: resourceJobStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"handle": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
