package main

import (
  "encoding/json"
  "fmt"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
)

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
  client := m.(PaperspaceClient).RestyClient

  log.Printf("[INFO] paperspace dataSourceUserRead Client ready")

  id := d.Get("id").(string)

  resp, err := client.R().
  Get("/users/getUsers?id=" + id)

  if err != nil {
    return fmt.Errorf("Error reading paperspace user: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace dataSourceUserRead StatusCode: %v", statusCode)
  LogResponse("paperspace dataSourceUserRead", resp, err)
  if statusCode == 404 {
    return fmt.Errorf("Error reading paperspace user: id not found %s",id)
  }
  if statusCode != 200 {
    return fmt.Errorf("Error reading paperspace user: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

  if err != nil {
    return fmt.Errorf("Error unmarshalling paperspace user read response: %s", err)
  }

  mpa := f.([]interface{})

  if len(mpa) > 1 {
    return fmt.Errorf("Error unmarshalling paperspace user read response: found more than one user for id %s", id)
  }
  if len(mpa) == 0 {
    return fmt.Errorf("Error unmarshalling paperspace user read response: user id not found %s", id)
  }

  mp, ok := mpa[0].(map[string]interface{})
  if !ok {
    return fmt.Errorf("Error unmarshalling paperspace user read response: user id not found %s", id)
  }

  idr, _ := mp["id"].(string)

  if idr == "" {
    return fmt.Errorf("Error unmarshalling paperspace user read response: user id not found %s", id)
  }

  if idr != id {
    return fmt.Errorf("Error unmarshalling paperspace user read response: found user id %s does not match id %v", idr, id)
  }

  log.Printf("[INFO] paperspace dataSourceUserRead user id: %v", idr)

  SetResData(d, mp, "email")
  SetResData(d, mp, "firstname")
  SetResData(d, mp, "lastname")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "teamId")

  d.SetId(idr)

	return nil
}

func dataSourceUser() *schema.Resource {
  return &schema.Resource{
    Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
      "id": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "email": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "firstname": &schema.Schema{
        Type:     schema.TypeString,
        Computed: true,
      },
      "lastname": &schema.Schema{
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
    },
	}
}
