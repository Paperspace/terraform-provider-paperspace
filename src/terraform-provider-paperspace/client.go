package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type MapIf map[string]interface{}

func (m *MapIf) Append(d *schema.ResourceData, k string) {
	v := d.Get(k)
	(*m)[k] = v
}

func (m *MapIf) AppendAs(d *schema.ResourceData, k, nk string) {
	v := d.Get(k)
	(*m)[nk] = v
}

func (m *MapIf) AppendV(d *schema.ResourceData, k, v string) {
	(*m)[k] = v
}

func (m *MapIf) AppendIfSet(d *schema.ResourceData, k string) {
	v := d.Get(k)
	if reflect.ValueOf(v).Interface() != reflect.Zero(reflect.TypeOf(v)).Interface() {
		(*m)[k] = v
	}
}

func (m *MapIf) AppendAsIfSet(d *schema.ResourceData, k, nk string) {
	v := d.Get(k)
	if reflect.ValueOf(v).Interface() != reflect.Zero(reflect.TypeOf(v)).Interface() {
		(*m)[nk] = v
	}
}

func SetResDataFrom(d *schema.ResourceData, m map[string]interface{}, dn, n string) {
	v, ok := m[n]
	//log.Printf("%v %v\n", n, v)
	if ok {
		d.Set(dn, v)
	}
}

func SetResData(d *schema.ResourceData, m map[string]interface{}, n string) {
	SetResDataFrom(d, m, n, n)
}

// LogObjectResponse logs http response fields
func LogObjectResponse(reqDesc string, reqURL *url.URL, resp *http.Response, body map[string]interface{}, err error) {
	log.Printf("[INFO] Request: %v", reqDesc)
	log.Printf("[INFO] Request URL: %v", reqURL)
	log.Printf("[INFO] Response Status: %v", resp.Status)
	log.Printf("[INFO] Response: %v", resp)
	log.Printf("[INFO] Response Body: %s", body)
	log.Printf("[INFO] Error: %v", err)
}

// LogArrayResponse logs http response fields
func LogArrayResponse(reqDesc string, reqURL *url.URL, resp *http.Response, body interface{}, err error) {
	log.Printf("[INFO] Request: %v", reqDesc)
	log.Printf("[INFO] Request URL: %v", reqURL)
	log.Printf("[INFO] Response Status: %v", resp.Status)
	log.Printf("[INFO] Response: %v", resp)
	log.Printf("[INFO] Response Body: %s", body)
	log.Printf("[INFO] Error: %v", err)
}

type ClientConfig struct {
	APIKey  string
	APIHost string
	Region  string
}

type PaperspaceClient struct {
	APIKey     string
	APIHost    string
	Region     string
	HttpClient *http.Client
}

func (c *ClientConfig) Client() (PaperspaceClient, error) {
	timeout := 10 * time.Second
	hc := &http.Client{
		Timeout: timeout,
	}

	rt := WithHeader(hc.Transport)
	rt.Set("x-api-key", c.APIKey)
	rt.Set("Accept", "application/json")
	rt.Set("Content-Type", "application/json")
	rt.Set("User-Agent", "terraform-provider-paperspace")
	rt.Set("ps_client_name", "terraform-provider-paperspace")
	hc.Transport = rt

	client := PaperspaceClient{
		APIKey:     c.APIKey,
		APIHost:    c.APIHost,
		Region:     c.Region,
		HttpClient: hc,
	}

	return client, nil
}

// from https://stackoverflow.com/questions/51325704/adding-a-default-http-header-in-go
type withHeader struct {
	http.Header
	rt http.RoundTripper
}

// WithHeader effectively allows http.Client to have global headers
func WithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{
		Header: make(http.Header),
		rt:     rt,
	}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}

func (psc *PaperspaceClient) GetMachine(id string) (body map[string]interface{}, err error) {
	url := fmt.Sprintf("%s/machines/getMachinePublic?machineId=%s", psc.APIHost, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error constructing GetMachine request: %s", err)
	}

	resp, err := psc.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending GetMachine request: %s", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("Error decoding GetMachine response body: %s", err)
	}

	LogObjectResponse("GetMachine", req.URL, resp, body, err)

	if resp.StatusCode != 404 && resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error on GetMachine response: statusCode: %d", resp.StatusCode)
	}

	nextID, _ := body["id"].(string)
	if resp.StatusCode == 404 || nextID == "" {
		return nil, fmt.Errorf("Error on GetMachine: machine not found")
	}

	return body, nil
}

func (psc *PaperspaceClient) CreateMachine(data []byte) (id string, err error) {
	url := fmt.Sprintf("%s/machines/createSingleMachinePublic", psc.APIHost)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("Error constructing CreateMachine request: %s", err)
	}

	resp, err := psc.HttpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending CreateMachine request: %s", err)
	}
	defer resp.Body.Close()

	body := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", fmt.Errorf("Error decoding CreateMachine response body: %s", err)
	}

	LogObjectResponse("CreateMachine", req.URL, resp, body, err)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error on CreateMachine: Response: %s", body)
	}

	id, _ = body["id"].(string)

	if id == "" {
		return "", fmt.Errorf("Error on CreateMachine: id not found")
	}

	return id, nil
}

func (psc *PaperspaceClient) DeleteMachine(id string) (err error) {
	url := fmt.Sprintf("%s/machines/%s/destroyMachine", psc.APIHost, id)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("Error constructing DeleteMachine request: %s", err)
	}

	resp, err := psc.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Error sending DeleteMachine request: %s", err)
	}
	defer resp.Body.Close()

	LogObjectResponse("DeleteMachine", req.URL, resp, nil, err)

	if resp.StatusCode != 204 {
		return fmt.Errorf("Error deleting machine, statusCode: %d", resp.StatusCode)
	}

	return nil
}
