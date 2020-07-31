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

func dataSourceUserRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := newInternalPaperspaceClient(m)

	log.Printf("[INFO] paperspace dataSourceUserRead Client ready")

	queryParam := false
	queryStr := "?"
	id, ok := d.GetOk("id")
	if ok {
		queryStr += "id=" + url.QueryEscape(id.(string))
		queryParam = true
	}
	email, ok := d.GetOk("email")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "email=" + url.QueryEscape(email.(string))
		queryParam = true
	}
	firstname, ok := d.GetOk("firstname")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "firstname=" + url.QueryEscape(firstname.(string))
		queryParam = true
	}
	lastname, ok := d.GetOk("lastname")
	if ok {
		if queryParam {
			queryStr += "&"
		}
		queryStr += "lastname=" + url.QueryEscape(lastname.(string))
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
	if !queryParam {
		return fmt.Errorf("Error reading paperspace user: must specify query filter properties")
	}

	url := fmt.Sprintf("%s/users/getUsers%s", paperspaceClient.APIHost, queryStr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error constructing GetUsers request: %s", err)
	}
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return fmt.Errorf("Error constructing GetUsers request: %s", err)
	}
	log.Print("[INFO] Request:", string(requestDump))

	resp, err := paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error reading paperspace user: %s", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	log.Printf("[INFO] paperspace dataSourceUserRead StatusCode: %v", statusCode)
	if statusCode == 404 {
		return fmt.Errorf("Error reading paperspace user: users not found")
	}
	if statusCode != 200 {
		responseDump, _ := httputil.DumpResponse(resp, true)
		return fmt.Errorf("Error reading paperspace user: Response: %s", string(responseDump))
	}

	var f interface{}
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return fmt.Errorf("Error decoding GetUsers response body: %s", err)
	}
	LogHttpResponse("paperspace dataSourceUserRead", req.URL, resp, f, err)

	mpa := f.([]interface{})
	if len(mpa) > 1 {
		return fmt.Errorf("Error reading paperspace user: found more than one user matching given properties")
	}
	if len(mpa) == 0 {
		return fmt.Errorf("Error reading paperspace user: no user found matching given properties")
	}

	mp, ok := mpa[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Error unmarshalling paperspace user read response: no users not found")
	}

	idr, _ := mp["id"].(string)
	if idr == "" {
		return fmt.Errorf("Error unmarshalling paperspace user read response: no user id found for user")
	}

	log.Printf("[INFO] paperspace dataSourceUserRead user id: %v", idr)

	SetResData(d, mp, "email")
	SetResData(d, mp, "firstname")
	SetResData(d, mp, "lastname")
	SetResDataFrom(d, mp, "dt_created", "dtCreated")
	SetResDataFrom(d, mp, "team_id", "teamId")

	d.SetId(idr)

	return nil
}

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": &schema.Schema{
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
			"dt_created": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
