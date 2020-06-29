package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func updateNetworkSchema(d *schema.ResourceData, network Network) {
	d.Set("network", network.Network)
	d.Set("netmask", network.Netmask)
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)
	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name is not a string")
	}
	teamID, ok := d.Get("team_id").(int)
	if !ok {
		return fmt.Errorf("team_id is not an int")
	}

	regionId, ok := RegionMap[paperspaceClient.Region]
	if !ok {
		return fmt.Errorf("Region %s not found", paperspaceClient.Region)
	}

	createNetworkParams := CreateNetworkParams{
		Name:     name,
		RegionId: regionId,
	}

	if err := paperspaceClient.CreateNetwork(teamID, createNetworkParams); err != nil {
		return fmt.Errorf("Error creating private network: %s", err)
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		paperspaceClient := m.(PaperspaceClient)

		// XXX: potential race condition for multiple networks created with the name concurrently
		// Add sync API response to API
		networks, err := paperspaceClient.GetTeamNetworks(teamID)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error creating private network: %s", err))
		}
		for _, network := range networks {
			if network.Handle == d.Id() {
				return resource.NonRetryableError(resourceNetworkRead(d, m))
			}
		}

		return resource.RetryableError(fmt.Errorf("Network not found"))
	})
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)
	teamID, ok := d.Get("team_id").(int)
	if !ok {
		return fmt.Errorf("team_id is not an int")
	}

	networks, err := paperspaceClient.GetTeamNetworks(teamID)
	if err != nil {
		return fmt.Errorf("Error creating private network: %s", err)
	}

	for _, network := range networks {
		if network.Handle == d.Id() {
			updateNetworkSchema(d, network)
			return nil
		}
	}

	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceNetworkRead(d, m)
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkCreate,
		Read:   resourceNetworkRead,
		Update: resourceNetworkUpdate,
		Delete: resourceNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}
