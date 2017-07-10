package main

import (
  "encoding/json"
  "fmt"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
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
  body.Append(d, "machineType")
  body.Append(d, "size")
  body.Append(d, "billingType")
  body.AppendAs(d, "name", "machineName")
  body.Append(d, "templateId")
  body.AppendIfSet(d, "networkId")
  body.AppendIfSet(d, "teamId")
  body.AppendIfSet(d, "userId")
  body.AppendIfSet(d, "email")
  body.AppendIfSet(d, "password")
  body.AppendIfSet(d, "firstName")
  body.AppendIfSet(d, "lastName")
  body.AppendIfSet(d, "notificationEmail")
  body.AppendIfSet(d, "scriptId")

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
    return fmt.Errorf("Error creating paperspace machine: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

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
  SetResData(d, mp, "storageTotal")
  SetResData(d, mp, "storageUsed")
  SetResData(d, mp, "usageRate")
  SetResData(d, mp, "shutdownTimeoutInHours")
  SetResData(d, mp, "shutdownTimeoutForces")
  SetResData(d, mp, "performAutoSnapshot")
  SetResData(d, mp, "autoSnapshotFrequency")
  SetResData(d, mp, "autoSnapshotSaveCount")
  SetResData(d, mp, "agentType")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "state")
  SetResData(d, mp, "networkId") //overlays with null initially
  SetResData(d, mp, "privateIpAddress")
  SetResData(d, mp, "publicIpAddress")
  SetResData(d, mp, "region") //overlays with null initially
  SetResData(d, mp, "userId")
  SetResData(d, mp, "teamId")
  SetResData(d, mp, "scriptId")
  SetResData(d, mp, "dtLastRun")

  d.SetId(id);

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
    return fmt.Errorf("Error reading paperspace machine: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

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
  SetResData(d, mp, "storageTotal")
  SetResData(d, mp, "storageUsed")
  SetResData(d, mp, "usageRate")
  SetResData(d, mp, "shutdownTimeoutInHours")
  SetResData(d, mp, "shutdownTimeoutForces")
  SetResData(d, mp, "performAutoSnapshot")
  SetResData(d, mp, "autoSnapshotFrequency")
  SetResData(d, mp, "autoSnapshotSaveCount")
  SetResData(d, mp, "agentType")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "state")
  SetResData(d, mp, "networkId") //overlays with null initially
  SetResData(d, mp, "privateIpAddress")
  SetResData(d, mp, "publicIpAddress")
  SetResData(d, mp, "region") //overlays with null initially
  SetResData(d, mp, "userId")
  SetResData(d, mp, "teamId")
  SetResData(d, mp, "scriptId")
  SetResData(d, mp, "dtLastRun")

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
    return fmt.Errorf("Error deleting paperspace machine: Response: %s", resp.Body())
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
      },
      "machineType": &schema.Schema{
          Type:     schema.TypeString,
          Required: true,
      },
      "size": &schema.Schema{
          Type:     schema.TypeInt,
          Required: true,
      },
      "billingType": &schema.Schema{
          Type:     schema.TypeString,
          Required: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "templateId": &schema.Schema{
          Type:     schema.TypeString,
          Required: true,
      },
      "networkId": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "teamId": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "userId": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "email": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "password": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "firstName": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "lastName": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "notificationEmail": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "scriptId": &schema.Schema{
          Type:     schema.TypeString,
          Optional: true,
      },
      "dtLastRun": &schema.Schema{
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
      "storageTotal": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "storageUsed": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "usageRate": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "shutdownTimeoutInHours": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "shutdownTimeoutForces": &schema.Schema{
          Type:     schema.TypeBool,
          Computed: true,
      },
      "performAutoSnapshot": &schema.Schema{
          Type:     schema.TypeBool,
          Computed: true,
      },
      "autoSnapshotFrequency": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "autoSnapshotSaveCount": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "agentType": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "dtCreated": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "state": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "privateIpAddress": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
      "publicIpAddress": &schema.Schema{
          Type:     schema.TypeString,
          Computed: true,
      },
    },
  }
}
