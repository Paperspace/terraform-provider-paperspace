package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var RegionMap = map[string]int{
	"East Coast (NY2)": 1,
	"West Coast (CA1)": 2,
	"Europe (AMS1)":    3,
}

type Network struct {
	ID      int    `json:"id"`
	Handle  string `json:"handle"`
	IsTaken bool   `json:"isTaken"`
	Network string `json:"network"`
	Netmask string `json:"netmask"`
	VlanID  int    `json:"vlanId"`
}

type NamedNetwork struct {
	Name    string  `json:"name"`
	Network Network `json:"network"`
}

type CreateTeamNamedNetworkParams struct {
	Name     string `json:"name"`
	RegionId int    `json:"regionId"`
}

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

func logHttpRequestConstruction(operationType string, url string, data *bytes.Buffer) {
	log.Printf("Constructing %s request to url: %s, data: %v", operationType, url, data)
}

// LogHttpResponse logs http response fields
func LogHttpResponse(reqDesc string, reqURL *url.URL, resp *http.Response, body interface{}, err error) {
	log.Printf("Request: %v", reqDesc)
	log.Printf("Request URL: %v", reqURL)
	log.Printf("Response Status: %v", resp.Status)
	log.Printf("Response: %v", resp)
	log.Printf("Response Body: %s", spew.Sdump(body))
	log.Printf("Error: %v", err)
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

func (c *ClientConfig) Client() (paperspaceClient PaperspaceClient, err error) {
	timeout := 10 * time.Second
	client := &http.Client{
		Timeout: timeout,
	}

	transport := WithHeader(client.Transport)
	transport.Set("x-api-key", c.APIKey)
	transport.Set("Accept", "application/json")
	transport.Set("Content-Type", "application/json")
	transport.Set("User-Agent", "terraform-provider-paperspace")
	transport.Set("ps_client_name", "terraform-provider-paperspace")
	client.Transport = transport

	paperspaceClient = PaperspaceClient{
		APIKey:     c.APIKey,
		APIHost:    c.APIHost,
		Region:     c.Region,
		HttpClient: client,
	}

	return paperspaceClient, nil
}

// from https://stackoverflow.com/questions/51325704/adding-a-default-http-header-in-go
type withHeader struct {
	http.Header
	transport http.RoundTripper
}

// WithHeader effectively allows http.Client to have global headers
func WithHeader(transport http.RoundTripper) withHeader {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return withHeader{
		Header:    make(http.Header),
		transport: transport,
	}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.transport.RoundTrip(req)
}

func (paperspaceClient *PaperspaceClient) RequestInterface(method string, url string, params, result interface{}) (res *http.Response, err error) {
	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, err
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return resp, err
	}

	LogHttpResponse("", req.URL, resp, result, err)
	return resp, nil
}

func (paperspaceClient *PaperspaceClient) Request(operationType string, url string, data []byte) (body map[string]interface{}, statusCode int, err error) {
	buf := bytes.NewBuffer(data)

	logHttpRequestConstruction(operationType, url, buf)

	req, err := http.NewRequest(operationType, url, buf)
	if err != nil {
		return nil, 0, fmt.Errorf("Error constructing request: %s", err)
	}

	resp, err := paperspaceClient.HttpClient.Do(req)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Error completing request: %s", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("Error decoding response body: %s", err)
	}

	LogHttpResponse("", req.URL, resp, body, err)

	return body, resp.StatusCode, nil
}

func (paperspaceClient *PaperspaceClient) GetMachine(id string) (body map[string]interface{}, err error) {
	url := fmt.Sprintf("%s/machines/getMachinePublic?machineId=%s", paperspaceClient.APIHost, id)
	body, statusCode, err := paperspaceClient.Request("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if statusCode != 404 && statusCode != 200 {
		return nil, fmt.Errorf("Error on GetMachine response: statusCode: %d", statusCode)
	}

	nextID, _ := body["id"].(string)
	if statusCode == 404 || nextID == "" {
		return nil, fmt.Errorf("Error on GetMachine: machine not found")
	}

	return body, nil
}

func (paperspaceClient *PaperspaceClient) CreateMachine(data []byte) (id string, err error) {
	url := fmt.Sprintf("%s/machines/createSingleMachinePublic", paperspaceClient.APIHost)
	body, statusCode, err := paperspaceClient.Request("POST", url, data)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("Error unmarshaling response body: %v", err)
		}

		return "", fmt.Errorf("Error on CreateMachine: Status Code %d, Response Body: %s", statusCode, jsonBody)
	}

	id, _ = body["id"].(string)

	if id == "" {
		return "", fmt.Errorf("Error on CreateMachine: id not found")
	}

	return id, nil
}

func (paperspaceClient *PaperspaceClient) DeleteMachine(id string) (err error) {
	url := fmt.Sprintf("%s/machines/%s/destroyMachine", paperspaceClient.APIHost, id)
	_, statusCode, err := paperspaceClient.Request("POST", url, nil)
	// /destroyMachine returns the string "EOF" if it was successful, which can't be JSON-decoded
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		return err
	}

	if statusCode != 204 {
		return fmt.Errorf("Error deleting machine")
	}

	return nil
}

func (paperspaceClient *PaperspaceClient) CreateTeamNamedNetwork(teamID int, createNamedNetworkParams CreateTeamNamedNetworkParams) error {
	var network Network
	url := fmt.Sprintf("%s/teams/%d/createPrivateNetwork", paperspaceClient.APIHost, teamID)

	_, err := paperspaceClient.RequestInterface("POST", url, createNamedNetworkParams, &network)
	if strings.Contains(err.Error(), "EOF") {
		return nil
	}
	return err
}

func (paperspaceClient *PaperspaceClient) GetTeamNamedNetworks(teamID int) ([]NamedNetwork, error) {
	var namedNetworks []NamedNetwork
	url := fmt.Sprintf("%s/teams/%d/getNetworks", paperspaceClient.APIHost, teamID)

	_, err := paperspaceClient.RequestInterface("GET", url, nil, &namedNetworks)
	spew.Sdump(namedNetworks)

	return namedNetworks, err
}

func (paperspaceClient *PaperspaceClient) GetTeamNamedNetwork(teamID int, name string) (*NamedNetwork, error) {
	namedNetworks, err := paperspaceClient.GetTeamNamedNetworks(teamID)
	if err != nil {
		return nil, err
	}

	for _, namedNetwork := range namedNetworks {
		if namedNetwork.Name == name {
			log.Print("FOUND NAMED NETWORK")
			spew.Sdump(namedNetwork)
			return &namedNetwork, nil
		}
	}

	return nil, fmt.Errorf("Error getting private network: %s", name)
}
