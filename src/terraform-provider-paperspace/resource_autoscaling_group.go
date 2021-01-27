package main

import (
	"time"

	"github.com/Paperspace/paperspace-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ErrNotFound(err error) bool {
	paperspaceError, ok := err.(*paperspace.PaperspaceError)
	if ok {
		if paperspaceError.Status == 404 {
			return true
		}
	}

	return false
}

func resourceAutoscalingGroupCreate(d *schema.ResourceData, m interface{}) error {
	var autoscalingGroup paperspace.AutoscalingGroup

	paperspaceClient := newPaperspaceClient(m)
	autoscalingGroupCreateParams := paperspace.AutoscalingGroupCreateParams{
		Name:        d.Get("name").(string),
		ClusterID:   d.Get("cluster_id").(string),
		Min:         d.Get("min").(int),
		Max:         d.Get("max").(int),
		MachineType: d.Get("machine_type").(string),
		TemplateID:  d.Get("template_id").(string),
		NetworkID:   d.Get("network_id").(string),
		ScriptID:    d.Get("startup_script_id").(string),
	}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		var err error
		autoscalingGroup, err = paperspaceClient.CreateAutoscalingGroup(autoscalingGroupCreateParams)
		if err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(nil)
	})
	if err != nil {
		return err
	}

	d.SetId(autoscalingGroup.ID)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if err := resourceAutoscalingGroupRead(d, m); err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(nil)
	})
}

func resourceAutoscalingGroupRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := newPaperspaceClient(m)

	autoscalingGroup, err := paperspaceClient.GetAutoscalingGroup(d.Id(), paperspace.AutoscalingGroupGetParams{})
	if err != nil {
		if ErrNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", autoscalingGroup.Name)
	d.Set("machine_type", autoscalingGroup.MachineType)
	d.Set("template_id", autoscalingGroup.TemplateID)
	d.Set("network_id", autoscalingGroup.NetworkID)
	d.Set("startup_script_id", autoscalingGroup.ScriptID)

	return nil
}

func resourceAutoscalingGroupUpdate(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := newPaperspaceClient(m)
	autoscalingGroupUpdateParams := paperspace.AutoscalingGroupUpdateParams{
		Attributes: paperspace.AutoscalingGroupUpdateAttributeParams{
			Name:        d.Get("name").(string),
			MachineType: d.Get("machine_type").(string),
			TemplateID:  d.Get("template_id").(string),
			NetworkID:   d.Get("network_id").(string),
			ScriptID:    d.Get("startup_script_id").(string),
		},
	}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if err := paperspaceClient.UpdateAutoscalingGroup(d.Id(), autoscalingGroupUpdateParams); err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(nil)
	})
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if err := resourceAutoscalingGroupRead(d, m); err != nil {
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(nil)
	})

}

func resourceAutoscalingGroupDelete(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := newPaperspaceClient(m)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		if err := paperspaceClient.DeleteAutoscalingGroup(d.Id(), paperspace.AutoscalingGroupDeleteParams{}); err != nil {
			if ErrNotFound(err) {
				return resource.NonRetryableError(nil)
			}
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(nil)
	})
}

func resourceAutoscalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutoscalingGroupCreate,
		Read:   resourceAutoscalingGroupRead,
		Update: resourceAutoscalingGroupUpdate,
		Delete: resourceAutoscalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"min": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"max": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"machine_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"template_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"startup_script_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}
