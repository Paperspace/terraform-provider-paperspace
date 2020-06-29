package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func generateNetworkHandle() string {
	rand.Seed(time.Now().UnixNano())

	return fmt.Sprint("ne" + randSeq(7))
}

func updateNetworkSchema(d *schema.ResourceData, network Network) {
	d.Set("network", network.Network)
	d.Set("netmask", network.Netmask)
}

func resourceNetworkCreate(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)
	teamID, ok := d.Get("team_id").(int)
	if !ok {
		return fmt.Errorf("team_id is not an int")
	}

	regionId, ok := RegionMap[paperspaceClient.Region]
	if !ok {
		return fmt.Errorf("Region %s not found", paperspaceClient.Region)
	}

	currentNetworks, err := paperspaceClient.GetTeamNetworks(teamID)
	if err != nil {
		return fmt.Errorf("Error getting current networks: %s", err)
	}
	spew.Sdump(currentNetworks)

	name := generateNetworkHandle()

	createNetworkParams := CreateNetworkParams{
		Name:     name,
		RegionId: regionId,
	}
	spew.Sdump(createNetworkParams)

	if err := paperspaceClient.CreateNetwork(teamID, createNetworkParams); err != nil {
		return fmt.Errorf("Error creating private network: %s", err)
	}
	d.SetId(name)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		paperspaceClient := m.(PaperspaceClient)

		// XXX: potential race condition for multiple networks created with the name concurrently
		// Add sync API response to API
		networks, err := paperspaceClient.GetTeamNetworks(teamID)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error creating private network: %s", err))
		}
		for _, network := range networks {
			if network.Name == d.Id() {
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
