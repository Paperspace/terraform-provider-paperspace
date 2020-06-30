package main

import (
	"fmt"
	"log"
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

func networkHandle() string {
	rand.Seed(time.Now().UnixNano())

	return fmt.Sprint("ne" + randSeq(7))
}

func updateNetworkSchema(d *schema.ResourceData, network Network, name string) {
	d.Set("handle", network.Handle)
	d.Set("is_taken", network.IsTaken)
	d.Set("name", name)
	d.Set("netmask", network.Netmask)
	d.Set("network", network.Network)
	d.Set("vlan_id", network.VlanID)
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

	name := networkHandle()

	createNamedNetworkParams := CreateTeamNamedNetworkParams{
		Name:     name,
		RegionId: regionId,
	}
	spew.Sdump(createNamedNetworkParams)

	if err := paperspaceClient.CreateTeamNamedNetwork(teamID, createNamedNetworkParams); err != nil {
		return fmt.Errorf("Error creating private network: %s", err)
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		paperspaceClient := m.(PaperspaceClient)

		// XXX: potential race condition for multiple networks created with the name concurrently
		// Add sync API response to API
		namedNetwork, err := paperspaceClient.GetTeamNamedNetwork(teamID, name)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error creating private network: %s", err))
		}

		d.SetId(string(namedNetwork.Network.ID))
		return resource.NonRetryableError(resourceNetworkRead(d, m))
	})
}

func resourceNetworkRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)
	teamID, ok := d.Get("team_id").(int)
	if !ok {
		return fmt.Errorf("team_id is not an int")
	}
	name, ok := d.Get("name").(string)
	if !ok {
		return fmt.Errorf("name is not a string")
	}

	namedNetwork, err := paperspaceClient.GetTeamNamedNetwork(teamID, name)
	if err != nil {
		d.SetId("")
		return err
	}
	log.Printf("!! bob")

	d.SetId(string(namedNetwork.Network.ID))
	updateNetworkSchema(d, namedNetwork.Network, name)

	spew.Sdump(d)
	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO: implement; api doesn't exist yet
	return resourceNetworkRead(d, m)
}

func resourceNetworkDelete(d *schema.ResourceData, m interface{}) error {
	// TODO: implement; api doesn't exist yet
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
			"team_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"handle": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_taken": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"netmask": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vlan_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			// name is not on the network schema but rather part of what i'm calling the
			// named network response, which comes from getNetworks. this is a joined
			// response between the network and network_owners table, to include name.
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
