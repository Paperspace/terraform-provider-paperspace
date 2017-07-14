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

  queryParam := false;
  queryStr := "?"
  id, ok := d.GetOk("id")
  if ok {
    queryStr += "id=" + id.(string)
    queryParam = true
  }
  name, ok := d.GetOk("name")
  if ok {
    if queryParam {
      queryStr += "&"
    }
    queryStr += "name=" + name.(string)
    queryParam = true
  }
  region, ok := d.GetOk("region")
  if ok {
    if queryParam {
      queryStr += "&"
    }
    queryStr += "region=" + region.(string)
    queryParam = true
  }
  dtCreated, ok := d.GetOk("dtCreated")
  if ok {
    if queryParam {
      queryStr += "&"
    }
    queryStr += "dtCreated=" + dtCreated.(string)
    queryParam = true
  }
  network, ok := d.GetOk("network")
  if ok {
    if queryParam {
      queryStr += "&"
    }
    queryStr += "network=" + network.(string)
    queryParam = true
  }
  netmask, ok := d.GetOk("netmask")
  if ok {
    if queryParam {
      queryStr += "&"
    }
    queryStr += "netmask=" + netmask.(string)
    queryParam = true
  }
  teamId, ok := d.GetOk("teamId")
  if ok {
    if queryParam {
      queryStr += "&"
    }
    queryStr += "teamId=" + teamId.(string)
    queryParam = true
  }
  if !queryParam {
    return fmt.Errorf("Error reading paperspace network: must specify query filter properties")
  }

  resp, err := client.R().
  Get("/networks/getNetworks" + queryStr)
  if err != nil {
    return fmt.Errorf("Error reading paperspace network: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace dataSourceNetworkRead StatusCode: %v", statusCode)
  LogResponse("paperspace dataSourceNetworkRead", resp, err)
  if statusCode == 404 {
    return fmt.Errorf("Error reading paperspace network: networks not found")
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
    return fmt.Errorf("Error reading paperspace network: found more than one network matching given properties")
  }
  if len(mpa) == 0 {
    return fmt.Errorf("Error reading paperspace network: no network found matching given properties")
  }

  mp, ok := mpa[0].(map[string]interface{})
  if !ok {
    return fmt.Errorf("Error unmarshalling paperspace network read response: no networks not found")
  }

  idr, _ := mp["id"].(string)
  if idr == "" {
    return fmt.Errorf("Error unmarshalling paperspace network read response: no network id found for network")
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
        Optional: true,
      },
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "label": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "os": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "dtCreated": &schema.Schema{
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
      "region": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
    },
	}
}
