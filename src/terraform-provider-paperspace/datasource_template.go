package main

import (
  "encoding/json"
  "fmt"
  "github.com/hashicorp/terraform/helper/schema"
  "log"
)

func dataSourceTemplateRead(d *schema.ResourceData, m interface{}) error {
  client := m.(PaperspaceClient).RestyClient

  log.Printf("[INFO] paperspace dataSourceTemplateRead Client ready")

  id := d.Get("id").(string)

  resp, err := client.R().
  Get("/templates/getTemplates?id=" + id)

  if err != nil {
    return fmt.Errorf("Error reading paperspace template: %s", err)
  }

  statusCode := resp.StatusCode()
  log.Printf("[INFO] paperspace dataSourceTemplateRead StatusCode: %v", statusCode)
  LogResponse("paperspace dataSourceTemplateRead", resp, err)
  if statusCode == 404 {
    return fmt.Errorf("Error reading paperspace template: id not found %s",id)
  }
  if statusCode != 200 {
    return fmt.Errorf("Error reading paperspace template: Response: %s", resp.Body())
  }

  var f interface{}
  err = json.Unmarshal(resp.Body(), &f)

  if err != nil {
    return fmt.Errorf("Error unmarshalling paperspace template read response: %s", err)
  }

  mpa := f.([]interface{})

  if len(mpa) > 1 {
    return fmt.Errorf("Error unmarshalling paperspace template read response: found more than one template for id %s", id)
  }
  if len(mpa) == 0 {
    return fmt.Errorf("Error unmarshalling paperspace template read response: template id not found %s", id)
  }

  mp, ok := mpa[0].(map[string]interface{})
  if !ok {
    return fmt.Errorf("Error unmarshalling paperspace template read response: template id not found %s", id)
  }

  idr, _ := mp["id"].(string)

  if idr == "" {
    return fmt.Errorf("Error unmarshalling paperspace template read response: template id not found %s", id)
  }

  if idr != id {
    return fmt.Errorf("Error unmarshalling paperspace template read response: found template id %s does not match id %v", idr, id)
  }

  log.Printf("[INFO] paperspace dataSourceTemplateRead template id: %v", idr)

  SetResData(d, mp, "name")
  SetResData(d, mp, "label")
  SetResData(d, mp, "os")
  SetResData(d, mp, "dtCreated")
  SetResData(d, mp, "teamId")
  SetResData(d, mp, "userId")
  SetResData(d, mp, "region")

  d.SetId(idr)

	return nil
}

func dataSourceTemplate() *schema.Resource {
  return &schema.Resource{
    Read: dataSourceTemplateRead,

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
