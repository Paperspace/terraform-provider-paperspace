package main

import (
  "encoding/json"
  "fmt"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
)

func dataSourceNetworkRead(d *schema.ResourceData, m interface{}) error {
  client := m.(PaperspaceClient).RestyClient

  log.Printf("[INFO] paperspace dataSourceNetworkRead Client ready")

  id := d.Get("id").(string)

  resp, err := client.R().
  Get("/networks/getNetworks?id=" + id)

  if err != nil {
    return fmt.Errorf("Error reading paperspace network: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace dataSourceNetworkRead StatusCode: %v", statusCode)
  LogResponse("paperspace dataSourceNetworkRead", resp, err)
  if statusCode == 404 {
    return fmt.Errorf("Error reading paperspace network: id not found %s",id)
  }
  if statusCode != 200 {
    return fmt.Errorf("Error reading paperspace network: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

  if err != nil {
    return fmt.Errorf("Error unmarshalling paperspace network read response: %s", err)
  }

  mpa := f.([]interface{})

  if len(mpa) > 1 {
    return fmt.Errorf("Error unmarshalling paperspace network read response: found more than one network for id %s", id)
  }
  if len(mpa) == 0 {
    return fmt.Errorf("Error unmarshalling paperspace network read response: network id not found %s", id)
  }

  mp, ok := mpa[0].(map[string]interface{})
  if !ok {
    return fmt.Errorf("Error unmarshalling paperspace network read response: network id not found %s", id)
  }

  idr, _ := mp["id"].(string)

  if idr == "" {
    return fmt.Errorf("Error unmarshalling paperspace network read response: network id not found %s", id)
  }

  if idr != id {
    return fmt.Errorf("Error unmarshalling paperspace network read response: found network id %s does not match id %v", idr, id)
  }

  log.Printf("[INFO] paperspace dataSourceNetworkRead network id: %v", idr)

  SetResData(d, mp, "name")
  SetResData(d, mp, "region")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "network")
  SetResData(d, mp, "netmask")
  SetResData(d, mp, "teamId")

  d.SetId(idr)

	return nil
}

func dataSourceNetwork() *schema.Resource {
  return &schema.Resource{
    Read: dataSourceNetworkRead,

		Schema: map[string]*schema.Schema{
      "id": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "label": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "os": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "dtCreated": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "teamId": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "userId": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "region": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
    },
	}
}
