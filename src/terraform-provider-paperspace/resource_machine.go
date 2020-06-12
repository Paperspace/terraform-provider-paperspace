package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	psc := m.(PaperspaceClient)

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
	body.AppendAsIfSet(d, "user_id", "userId")
	body.AppendAsIfSet(d, "team_id", "teamId")
	body.AppendAsIfSet(d, "script_id", "scriptId")

	// fields not tested when this project was picked back up for https://github.com/Paperspace/terraform-provider-paperspace/pull/3
	body.AppendAsIfSet(d, "network_id", "networkId")
	body.AppendIfSet(d, "email")
	body.AppendIfSet(d, "password")
	body.AppendAsIfSet(d, "firstname", "firstName")
	body.AppendAsIfSet(d, "lastname", "lastName")
	body.AppendAsIfSet(d, "notification_email", "notificationEmail")

	data, _ := json.MarshalIndent(body, "", "  ")
	log.Println(string(data))

	id, err := psc.CreateMachine(data)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		body, _, err := psc.GetMachine(*id)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		state, ok := body["state"].(string)
		if !ok {
			return resource.RetryableError(fmt.Errorf("[WARNING] Expected machine to be ready but found no state"))
		}
		if state != "ready" {
			return resource.RetryableError(fmt.Errorf("[INFO] Expected machine to be ready but was in state %s", state))
		}

		SetResData(d, body, "name")
		SetResData(d, body, "os")
		SetResData(d, body, "ram")
		SetResData(d, body, "cpus")
		SetResData(d, body, "gpu")
		SetResDataFrom(d, body, "storage_total", "storageTotal")
		SetResDataFrom(d, body, "storage_used", "storageUsed")
		SetResDataFrom(d, body, "usage_rate", "usageRate")
		SetResDataFrom(d, body, "shutdown_timeout_in_hours", "shutdownTimeoutInHours")
		SetResDataFrom(d, body, "shutdown_timeout_forces", "shutdownTimeoutForces")
		SetResDataFrom(d, body, "perform_auto_snapshot", "performAutoSnapshot")
		SetResDataFrom(d, body, "auto_snapshot_frequency", "autoSnapshotFrequency")
		SetResDataFrom(d, body, "auto_snapshot_save_count", "autoSnapshotSaveCount")
		SetResDataFrom(d, body, "agent_type", "agentType")
		SetResDataFrom(d, body, "dt_created", "dtCreated")
		SetResData(d, body, "state")
		SetResDataFrom(d, body, "network_id", "networkId") //overlays with null initially
		SetResDataFrom(d, body, "private_ip_address", "privateIpAddress")
		SetResDataFrom(d, body, "public_ip_address", "publicIpAddress")
		SetResData(d, body, "region") //overlays with null initially
		SetResDataFrom(d, body, "user_id", "userId")
		SetResDataFrom(d, body, "team_id", "teamId")
		SetResDataFrom(d, body, "script_id", "scriptId")
		SetResDataFrom(d, body, "dt_last_run", "dtLastRun")

		d.SetId(*id)

		return resource.NonRetryableError(resourceMachineRead(d, m))
	})
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	psc := m.(PaperspaceClient)

	body, statusCode, err := psc.GetMachine(d.Id())
	if err != nil {
		return err
	}

	if *statusCode == 404 {
		log.Printf("[INFO] paperspace resourceMachineRead machineId not found; removing resource %s", d.Id())
		d.SetId("")
		return nil
	}
	if *statusCode != 200 {
		return fmt.Errorf("Error reading paperspace machine: Response: %s", body)
	}

	id, _ := body["id"].(string)

	if id == "" {
		log.Printf("[WARNING] paperspace resourceMachineRead machine id not found; removing resource %s", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] paperspace resourceMachineRead returned id: %v", id)

	mp := make(map[string]interface{})
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

	return resourceMachineRead(d, m)
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(PaperspaceClient).HttpClient

	log.Printf("[INFO] paperspace resourceMachineDelete Client ready")

	url := fmt.Sprintf("/machines/%s/destroyMachine", d.Id())
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Errorf("[WARNING] Constructing resourceMachineCreate delete machine request failed: %s", err)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error deleting paperspace machine: %s", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("Error deleting paperspace machine: Response: %s", resp.Body)
	}

	log.Printf("[INFO] paperspace resourceMachineDelete machine successfully started deleting, StatusCode: %v", resp.StatusCode)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		url := fmt.Sprintf("%s/machines/getMachinePublic?machineId=%s", m.(PaperspaceClient).APIHost, d.Id())
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[WARNING] Constructing resourceMachineCreate get machine request failed: %s", err))
		}

		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error getting paperspace machine: %s", err))
		}
		log.Printf("[DEBUG] paperspace resourceMachineCreate get machine response StatusCode: %v", resp.StatusCode)

		mp := make(map[string]interface{})
		err = json.NewDecoder(resp.Body).Decode(&mp)

		LogResponse("paperspace resourceMachineDelete", resp, err)

		if resp.StatusCode != 200 && resp.StatusCode != 404 {
			return resource.NonRetryableError(fmt.Errorf("Error getting paperspace machine, Status Code %s", err))
		}

		if resp.StatusCode == 200 {
			state, _ := mp["state"].(string)
			return resource.RetryableError(fmt.Errorf("Expected machine to be deleted but was in state %s", state))
		}

		// resp.StatusCode == 404
		log.Printf("[INFO] paperspace resourceMachineDelete machine successfully deleted, StatusCode: %v", resp.StatusCode)
		return resource.NonRetryableError(nil)
	})
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
				Type:     schema.TypeInt,
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}
