package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceScriptCreate(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

	log.Printf("[INFO] paperspace resourceScriptCreate Client ready")

	region := m.(PaperspaceClient).Region
	if r, ok := d.GetOk("region"); ok {
		region = r.(string)
	}
	if region == "" {
		return fmt.Errorf("Error creating paperspace script: missing region")
	}

	body := make(MapIf)
	body.AppendAs(d, "name", "scriptName")
	body.AppendAs(d, "script_text", "scriptText")
	body.AppendAsIfSet(d, "description", "scriptDescription")
	body.AppendAsIfSet(d, "is_enabled", "isEnabled")
	body.AppendAsIfSet(d, "run_once", "runOnce")

	data, _ := json.MarshalIndent(body, "", "  ")
	log.Println(string(data))

	url := fmt.Sprintf("%s/scripts/createScript", paperspaceClient.APIHost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Error constructing CreateScript request: %s", err)
	}

	resp, err := paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error creating paperspace script: %s", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	log.Printf("[INFO] paperspace resourceScriptCreate StatusCode: %v", statusCode)
	if statusCode != 200 {
		return fmt.Errorf("Error creating paperspace script: Response: %s", resp.Body)
	}

	var f interface{}
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return fmt.Errorf("Error decoding GetScript response body: %s", err)
	}
	LogHttpResponseArray("paperspace dataSourceGetScript", req.URL, resp, f, err)

	mp := f.(map[string]interface{})
	id, _ := mp["id"].(string)

	if id == "" {
		return fmt.Errorf("Error in paperspace script create data: id not found")
	}

	log.Printf("[INFO] paperspace resourceScriptCreate returned id: %v", id)

	SetResData(d, mp, "name")
	SetResData(d, mp, "description")
	SetResDataFrom(d, mp, "owner_type", "ownerType")
	SetResDataFrom(d, mp, "owner_id", "ownerId")
	SetResDataFrom(d, mp, "dt_created", "dtCreated")
	SetResDataFrom(d, mp, "is_enabled", "isEnabled")
	SetResDataFrom(d, mp, "run_once", "runOnce")

	d.SetId(id)

	return resourceScriptRead(d, m)
}

func resourceScriptRead(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

	log.Printf("[INFO] paperspace resourceScriptRead Client ready")

	url := fmt.Sprintf("%s/scripts/getScript?scriptId=%s", paperspaceClient.APIHost, d.Id())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error constructing GetScript request: %s", err)
	}

	resp, err := paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error completing GetScript request: %s", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	log.Printf("[INFO] paperspace resourceScriptRead StatusCode: %v", statusCode)
	if statusCode == 404 {
		log.Printf("[INFO] paperspace resourceScriptRead scriptId not found; removing resource %s", d.Id())
		d.SetId("")
		return nil
	}
	var body interface{}
	json.NewDecoder(resp.Body).Decode(&body)
	LogHttpResponseArray("paperspace resourceScriptCreate", req.URL, resp, body, err)

	if statusCode != 200 {
		return fmt.Errorf("Error reading paperspace script: Response: %s", body)
	}

	if err != nil {
		return fmt.Errorf("Error unmarshalling paperspace script read response: %s", err)
	}

	mp := body.(map[string]interface{})
	id, _ := mp["id"].(string)

	if id == "" {
		log.Printf("[WARNING] paperspace resourceScriptRead script id not found; removing resource %s", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] paperspace resourceScriptRead returned id: %v", id)

	SetResData(d, mp, "name")
	SetResData(d, mp, "description")
	SetResDataFrom(d, mp, "owner_type", "ownerType")
	SetResDataFrom(d, mp, "owner_id", "ownerId")
	SetResDataFrom(d, mp, "dt_created", "dtCreated")
	SetResDataFrom(d, mp, "is_enabled", "isEnabled")
	SetResDataFrom(d, mp, "run_once", "runOnce")

	url = fmt.Sprintf("%s/scripts/getScriptText?scriptId=%s", paperspaceClient.APIHost, d.Id())
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error constructing GetScriptText request: %s", err)
	}

	resp, err = paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error reading paperspace script text: %s", err)
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	log.Printf("[INFO] paperspace resourceScriptRead text StatusCode: %v", statusCode)

	json.NewDecoder(resp.Body).Decode(&body)
	s, err := json.Marshal(body)
	LogHttpResponseArray("paperspace resourceScriptCreate", req.URL, resp, body, err)

	if statusCode == 404 {
		log.Printf("[INFO] paperspace resourceScriptRead text scriptId not found")
		return nil
	}
	if statusCode != 200 {
		return fmt.Errorf("Error reading paperspace script text: Response: %s", body)
	}

	d.Set("script_text", s)

	return nil
}

func resourceScriptUpdate(d *schema.ResourceData, m interface{}) error {

	log.Printf("[INFO] paperspace resourceScriptUpdate Client ready")

	return resourceScriptRead(d, m)
}

func resourceScriptDelete(d *schema.ResourceData, m interface{}) error {
	paperspaceClient := m.(PaperspaceClient)

	log.Printf("[INFO] paperspace resourceScriptDelete Client ready")

	url := fmt.Sprintf("%s/scripts/%s/destroy", paperspaceClient.APIHost, d.Id())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("Error constructing DeleteScript request: %s", err)
	}
	resp, err := paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error deleting paperspace script: %s", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	log.Printf("[INFO] paperspace resourceScriptDelete StatusCode: %v", statusCode)
	LogHttpResponseArray("paperspace resourceScriptDelete", req.URL, resp, nil, err)
	if statusCode != 204 && statusCode != 404 {
		return fmt.Errorf("Error deleting paperspace script: Response: %s", resp.Body)
	}
	if statusCode == 204 {
		log.Printf("[INFO] paperspace resourceScriptDelete script deleted successfully, StatusCode: %v", statusCode)
	}
	if statusCode == 404 {
		log.Printf("[INFO] paperspace resourceScriptDelete script already deleted, StatusCode: %v", statusCode)
	}

	return nil
}

func resourceScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceScriptCreate,
		Read:   resourceScriptRead,
		Update: resourceScriptUpdate,
		Delete: resourceScriptDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"script_text": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"owner_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"dt_created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"run_once": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}
