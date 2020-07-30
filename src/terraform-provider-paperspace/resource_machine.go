package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceMachineCreate(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

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
	body.AppendAsIfSet(d, "network_id", "networkId")
	body.AppendAsIfSet(d, "shutdown_timeout_in_hours", "shutdownTimeoutInHours")
	body.AppendAsIfSet(d, "is_managed", "isManaged")

	s := d.Get("live_forever")
	if s.(bool) == true {
		body["shutdownTimeoutInHours"] = nil
	}

	// fields not tested when this project was picked back up for https://github.com/Paperspace/terraform-provider-paperspace/pull/3
	body.AppendIfSet(d, "email")
	body.AppendIfSet(d, "password")
	body.AppendAsIfSet(d, "firstname", "firstName")
	body.AppendAsIfSet(d, "lastname", "lastName")
	body.AppendAsIfSet(d, "notification_email", "notificationEmail")

	data, _ := json.MarshalIndent(body, "", "  ")

	id, err := paperspaceClient.CreateMachine(data)
	if err != nil {
		return err
	}
	d.SetId(id)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		body, err := paperspaceClient.GetMachine(id)
		if err != nil {
			return resource.RetryableError(err)
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
		SetResDataFrom(d, body, "is_managed", "isManaged")

		return resource.NonRetryableError(resourceMachineRead(d, m))
	})
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

	_, err := paperspaceClient.GetMachine(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}

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
	SetResDataFrom(d, mp, "is_managed", "isManaged")

	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	// TODO: needs to be implemented
	return resourceMachineRead(d, m)
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

	err := paperspaceClient.DeleteMachine(d.Id())
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		body, err := paperspaceClient.GetMachine(d.Id())
		log.Printf("\nbody: %v\nerr: %v", body, err)
		if err != nil {
			if strings.Contains(err.Error(), "machine not found") {
				return resource.NonRetryableError(nil)
			}
			return resource.RetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Expected machine to be deleted but still exists"))
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
				Optional: true,
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
			"live_forever": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_managed": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}
