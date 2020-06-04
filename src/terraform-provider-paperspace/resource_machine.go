package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(PaperspaceClient).RestyClient

	log.Printf("[INFO] paperspace resourceMachineCreate Client ready")

	region := m.(PaperspaceClient).Region
	if r, ok := d.GetOk("region"); ok {
		region = r.(string)
	}
	if region == "" {
		return fmt.Errorf("Error creating paperspace machine: missing region")
	}

	body := make(MapIf)
	body.AppendV(d, "region", region)
	body.AppendAs(d, "machine_type", "machineType")
	body.Append(d, "size")
	body.AppendAs(d, "billing_type", "billingType")
	body.AppendAs(d, "name", "machineName")
	body.AppendAs(d, "template_id", "templateId")
	body.AppendAsIfSet(d, "assign_public_ip", "assignPublicIp")
	body.AppendAsIfSet(d, "network_id", "networkId")
	body.AppendAsIfSet(d, "team_id", "teamId")
	body.AppendAsIfSet(d, "user_id", "userId")
	body.AppendIfSet(d, "email")
	body.AppendIfSet(d, "password")
	body.AppendAsIfSet(d, "firstname", "firstName")
	body.AppendAsIfSet(d, "lastname", "lastName")
	body.AppendAsIfSet(d, "notification_email", "notificationEmail")
	body.AppendAsIfSet(d, "script_id", "scriptId")

	data, _ := json.MarshalIndent(body, "", "  ")
	log.Println(string(data))

	resp, err := client.R().
		SetBody(body).
		Post("/machines/createSingleMachinePublic")

	if err != nil {
		return fmt.Errorf("Error creating paperspace machine: %s", err)
	}

	statusCode := resp.StatusCode()
	log.Printf("[INFO] paperspace resourceMachineCreate StatusCode: %v", statusCode)
	LogResponse("paperspace resourceMachineCreate", resp, err)
	if statusCode != 200 {
		return fmt.Errorf("Error creating paperspace machine: Response: %s", resp.Body)
	}

	var f interface{}
	err = json.Unmarshal(resp.Body, &f)

	/*fake := []byte(`{"id":"psmfffm3","name":"Tom Terraform Test 4","os":null,"ram":null,
	  "cpus":1,"gpu":null,"storageTotal":null,"storageUsed":null,"usageRate":"C1 Hourly",
	  "shutdownTimeoutInHours":null,"shutdownTimeoutForces":false,"performAutoSnapshot":false,
	  "autoSnapshotFrequency":null,"autoSnapshotSaveCount":null,"agentType":"LinuxHeadless",
	  "dtCreated":"2017-06-22T04:29:59.501Z","state":"provisioning","networkId":null,
	  "privateIpAddress":null,"publicIpAddress":null,"region":null,"userId":"uijn3il",
	  "teamId":null}`)
	err := json.Unmarshal(fake, &f)*/

	if err != nil {
		return fmt.Errorf("Error unmarshalling paperspace machine create response: %s", err)
	}

	mp := f.(map[string]interface{})
	id, _ := mp["id"].(string)

	if id == "" {
		return fmt.Errorf("Error in paperspace machine create data: id not found")
	}

	log.Printf("[INFO] paperspace resourceMachineCreate returned id: %v", id)

	SetResData(d, mp, "name")
	SetResData(d, mp, "os")
	SetResData(d, mp, "ram")
	SetResData(d, mp, "cpus")
	SetResData(d, mp, "gpu")
	SetResDataFrom(d, mp, "storage_total", "storageTotal")
	SetResDataFrom(d, mp, "storage_used", "storageUsed")
	SetResDataFrom(d, mp, "usage_rate", "usageRate")
	SetResDataFrom(d, mp, "shutdown_timeout_in_hours", "shutdownTimeoutInHours")
	SetResDataFrom(d, mp, "shutdown_timeout_forces", "shutdownTimeoutForces")
	SetResDataFrom(d, mp, "perform_auto_snapshot", "performAutoSnapshot")
	SetResDataFrom(d, mp, "auto_snapshot_frequency", "autoSnapshotFrequency")
	SetResDataFrom(d, mp, "auto_snapshot_save_count", "autoSnapshotSaveCount")
	SetResDataFrom(d, mp, "agent_type", "agentType")
	SetResDataFrom(d, mp, "dt_created", "dtCreated")
	SetResData(d, mp, "state")
	SetResDataFrom(d, mp, "network_id", "networkId") //overlays with null initially
	SetResDataFrom(d, mp, "private_ip_address", "privateIpAddress")
	SetResDataFrom(d, mp, "public_ip_address", "publicIpAddress")
	SetResData(d, mp, "region") //overlays with null initially
	SetResDataFrom(d, mp, "user_id", "userId")
	SetResDataFrom(d, mp, "team_id", "teamId")
	SetResDataFrom(d, mp, "script_id", "scriptId")
	SetResDataFrom(d, mp, "dt_last_run", "dtLastRun")

	d.SetId(id)

	return nil
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	client := m.(PaperspaceClient).RestyClient

	log.Printf("[INFO] paperspace resourceMachineRead Client ready")

	resp, err := client.R().
		Get("/machines/getMachinePublic?machineId=" + d.Id())

	if err != nil {
		return fmt.Errorf("Error reading paperspace machine: %s", err)
	}

	statusCode := resp.StatusCode()
	log.Printf("[INFO] paperspace resourceMachineRead StatusCode: %v", statusCode)
	LogResponse("paperspace resourceMachineCreate", resp, err)
	if statusCode == 404 {
		log.Printf("[INFO] paperspace resourceMachineRead machineId not found; removing resource %s", d.Id())
		d.SetId("")
		return nil
	}
	if statusCode != 200 {
		return fmt.Errorf("Error reading paperspace machine: Response: %s", resp.Body)
	}

	var f interface{}
	err = json.Unmarshal(resp.Body, &f)

	if err != nil {
		return fmt.Errorf("Error unmarshalling paperspace machine read response: %s", err)
	}

	mp := f.(map[string]interface{})
	id, _ := mp["id"].(string)

	if id == "" {
		log.Printf("[WARNING] paperspace resourceMachineRead machine id not found; removing resource %s", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] paperspace resourceMachineRead returned id: %v", id)

	SetResData(d, mp, "name")
	SetResData(d, mp, "os")
	SetResData(d, mp, "ram")
	SetResData(d, mp, "cpus")
	SetResData(d, mp, "gpu")
	SetResDataFrom(d, mp, "storage_total", "storageTotal")
	SetResDataFrom(d, mp, "storage_used", "storageUsed")
	SetResDataFrom(d, mp, "usage_rate", "usageRate")
	SetResDataFrom(d, mp, "shutdown_timeout_in_hours", "shutdownTimeoutInHours")
	SetResDataFrom(d, mp, "shutdown_timeout_forces", "shutdownTimeoutForces")
	SetResDataFrom(d, mp, "perform_auto_snapshot", "performAutoSnapshot")
	SetResDataFrom(d, mp, "auto_snapshot_frequency", "autoSnapshotFrequency")
	SetResDataFrom(d, mp, "auto_snapshot_save_count", "autoSnapshotSaveCount")
	SetResDataFrom(d, mp, "agent_type", "agentType")
	SetResDataFrom(d, mp, "dt_created", "dtCreated")
	SetResData(d, mp, "state")
	SetResDataFrom(d, mp, "network_id", "networkId") //overlays with null initially
	SetResDataFrom(d, mp, "private_ip_address", "privateIpAddress")
	SetResDataFrom(d, mp, "public_ip_address", "publicIpAddress")
	SetResData(d, mp, "region") //overlays with null initially
	SetResDataFrom(d, mp, "user_id", "userId")
	SetResDataFrom(d, mp, "team_id", "teamId")
	SetResDataFrom(d, mp, "script_id", "scriptId")
	SetResDataFrom(d, mp, "dt_last_run", "dtLastRun")

	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {

	log.Printf("[INFO] paperspace resourceMachineUpdate Client ready")

	return nil
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(PaperspaceClient).RestyClient

	log.Printf("[INFO] paperspace resourceMachineDelete Client ready")

	resp, err := client.R().
		Post("/machines/" + d.Id() + "/destroyMachine")

	if err != nil {
		return fmt.Errorf("Error deleting paperspace machine: %s", err)
	}

	statusCode := resp.StatusCode()
	log.Printf("[INFO] paperspace resourceMachineDelete StatusCode: %v", statusCode)
	LogResponse("paperspace resourceMachineDelete", resp, err)
	if statusCode != 204 && statusCode != 404 {
		return fmt.Errorf("Error deleting paperspace machine: Response: %s", resp.Body)
	}
	if statusCode == 204 {
		log.Printf("[INFO] paperspace resourceMachineDelete machine deleted successfully, StatusCode: %v", statusCode)
	}
	if statusCode == 404 {
		log.Printf("[INFO] paperspace resourceMachineDelete machine already deleted, StatusCode: %v", statusCode)
	}

	return nil
}

func resourceMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMachineCreate,
		Read:   resourceMachineRead,
		Update: resourceMachineUpdate,
		Delete: resourceMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"machine_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"billing_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"assign_public_ip": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"firstname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"script_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"dt_last_run": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"os": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gpu": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_total": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_used": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage_rate": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"shutdown_timeout_in_hours": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"shutdown_timeout_forces": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"perform_auto_snapshot": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_snapshot_frequency": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_snapshot_save_count": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dt_created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
