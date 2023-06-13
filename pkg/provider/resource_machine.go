package provider

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
	paperspaceClient := newInternalPaperspaceClient(m)

	region := paperspaceClient.Region
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
	body.AppendAsIfSet(d, "perform_auto_snapshot", "performAutoSnapshot")
	body.AppendAsIfSet(d, "auto_snapshot_frequency", "autoSnapshotFrequency")
	body.AppendAsIfSet(d, "auto_snapshot_save_count", "autoSnapshotSaveCount")

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

		return resource.NonRetryableError(resourceMachineRead(d, m))
	})
}

func resourceMachineRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := newInternalPaperspaceClient(m)

	body, err := paperspaceClient.GetMachine(d.Id())
	if err != nil {
		if err.Error() == MachineNotFoundError {
			d.SetId("")
			return nil
		}

		return err
	}

	SetResData(d, body, "name")
	SetResData(d, body, "os")
	SetResData(d, body, "ram")
	SetResData(d, body, "cpus")
	SetResData(d, body, "gpu")
	SetResDataFrom(d, body, "storage_total", "storageTotal")
	SetResDataFrom(d, body, "storage_used", "storageUsed")
	SetResDataFrom(d, body, "usage_rate", "usageRate")

	shutdown_timeout := d.Get("shutdown_timeout_in_hours")
	_, ok := shutdown_timeout.(int32)
	if ok {
		SetResDataFrom(d, body, "shutdown_timeout_in_hours", "shutdownTimeoutInHours")
	}

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

	return nil
}

func resourceMachineUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceMachineRead(d, m)
}

func resourceMachineDelete(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := newInternalPaperspaceClient(m)

	err := paperspaceClient.DeleteMachine(d.Id())
	if err != nil {
		if err.Error() == MachineDeleteNotFoundError {
			d.SetId("")
			return nil
		}
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
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"billing_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"assign_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"script_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dt_last_run": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"gpu": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_total": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_used": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"usage_rate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shutdown_timeout_in_hours": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"shutdown_timeout_forces": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"perform_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_snapshot_frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_snapshot_save_count": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"agent_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dt_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"live_forever": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_managed": {
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
