package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTemplateRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

	log.Printf("[INFO] paperspace dataSourceTemplateRead Client ready")

	queryParam := false
	queryStr := "?"
	id, ok := d.GetOk("id")
	if ok {
		queryStr += "id=" + url.QueryEscape(id.(string))
		queryParam = true
	}
	name, ok := d.GetOk("name")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "name=" + url.QueryEscape(name.(string))
		queryParam = true
	}
	label, ok := d.GetOk("label")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "label=" + url.QueryEscape(label.(string))
		queryParam = true
	}
	os, ok := d.GetOk("os")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "os=" + url.QueryEscape(os.(string))
		queryParam = true
	}
	dtCreated, ok := d.GetOk("dt_created")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "dtCreated=" + url.QueryEscape(dtCreated.(string))
		queryParam = true
	}
	teamId, ok := d.GetOk("team_id")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "teamId=" + url.QueryEscape(teamId.(string))
		queryParam = true
	}
	userId, ok := d.GetOk("user_id")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "userId=" + url.QueryEscape(userId.(string))
		queryParam = true
	}
	region, ok := d.GetOk("region")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "region=" + url.QueryEscape(region.(string))
		queryParam = true
	}
	if !queryParam {
		return fmt.Errorf("Error reading paperspace template: must specify query filter properties")
	}

	url := fmt.Sprintf("%s/templates/getTemplates%s", paperspaceClient.APIHost, queryStr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error constructing GetTemplates request: %s", err)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return fmt.Errorf("Error constructing GetTemplates request: %s", err)
	}
	log.Print("[INFO] Request:", string(requestDump))

	resp, err := paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error reading paperspace template: %s", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	log.Printf("[INFO] paperspace dataSourceTemplateRead StatusCode: %v", statusCode)
	if statusCode == 404 {
		return fmt.Errorf("Error reading paperspace template: templates not found")
	}
	if statusCode != 200 {
		responseDump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("Error reading paperspace template: Response: %s", string(responseDump))
	}

	var f interface{}
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return fmt.Errorf("Error decoding GetTemplate response body: %s", err)
	}
	LogHttpResponse("paperspace dataSourceTemplateRead", req.URL, resp, f, err)

	mpa := f.([]interface{})
	if len(mpa) > 1 {
		return fmt.Errorf("Error reading paperspace template: found more than one template matching given properties")
	}
	if len(mpa) == 0 {
		return fmt.Errorf("Error reading paperspace template: no template found matching given properties")
	}

	mp, ok := mpa[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error unmarshalling paperspace template read response: no templates not found")
	}

	idr, _ := mp["id"].(string)
	if idr == "" {
		return fmt.Errorf("Error unmarshalling paperspace template read response: no template id found for template")
	}

	log.Printf("[INFO] paperspace dataSourceTemplateRead template id: %v", idr)

	SetResData(d, mp, "name")
	SetResData(d, mp, "label")
	SetResData(d, mp, "os")
	SetResDataFrom(d, mp, "dt_created", "dtCreated")
	SetResDataFrom(d, mp, "team_id", "teamId")
	SetResDataFrom(d, mp, "user_id", "userId")
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
			"dt_created": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": &schema.Schema{
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
