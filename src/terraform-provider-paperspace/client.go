package main

import (
	"log"
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/resty.v0"
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

type Config struct {
	ApiKey  string
	ApiHost string
	Region  string
}

type PaperspaceClient struct {
	ApiKey      string
	ApiHost     string
	Region      string
	RestyClient *resty.Client
}

func (c *Config) Client() (PaperspaceClient, error) {

	restyClient := resty.New()

	restyClient.
		SetHostURL(c.ApiHost).
		SetHeader("x-api-key", c.ApiKey).
		SetHeader("Accept", "application/json")

	client := PaperspaceClient{
		ApiKey:      c.ApiKey,
		ApiHost:     c.ApiHost,
		Region:      c.Region,
		RestyClient: restyClient,
	}

	return client, nil
}

func LogResponse(reqDesc string, resp *resty.Response, err error) {
	log.Printf("Request: %v", reqDesc)
	log.Printf("Error: %v", err)
	log.Printf("Response Status Code: %v", resp.StatusCode())
	log.Printf("Response Status: %v", resp.Status())
	log.Printf("Response Time: %v", resp.Time())
	log.Printf("Response Recevied At: %v", resp.ReceivedAt())
	log.Printf("Response Body: %v", resp) // or resp.String() or string(resp.Body)
}
