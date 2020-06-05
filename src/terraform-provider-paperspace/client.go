package main

import (
	"log"
	"reflect"

	"github.com/go-resty/resty/v2"
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

	restyClient.SetDebug(true)

	restyClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Explore request object
		log.Println("Request Info:")
		log.Println("client.HostURL", c.HostURL)
		log.Println("client.Header", c.Header)
		log.Println("req.Method", req.Method)
		log.Println("req.URL", req.URL)
		log.Println("req.Body", req.Body)
		log.Println("req.AuthScheme", req.AuthScheme)
		log.Println("req.RawRequest", req.RawRequest)
		log.Println("req.Error", req.Error)

		return nil
	})

	restyClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {

		// Explore response object
		log.Println("Response Info:")
		// log.Println("Error      :", err)
		log.Println("Status Code:", resp.StatusCode())
		log.Println("Status     :", resp.Status())
		log.Println("Proto      :", resp.Proto())
		log.Println("Time       :", resp.Time())
		log.Println("Received At:", resp.ReceivedAt())
		log.Println("Body       :\n", resp)
		log.Println()

		// Explore trace info
		log.Println("Request Trace Info:")
		ti := resp.Request.TraceInfo()
		log.Println("DNSLookup    :", ti.DNSLookup)
		log.Println("ConnTime     :", ti.ConnTime)
		log.Println("TCPConnTime  :", ti.TCPConnTime)
		log.Println("TLSHandshake :", ti.TLSHandshake)
		log.Println("ServerTime   :", ti.ServerTime)
		log.Println("ResponseTime :", ti.ResponseTime)
		log.Println("TotalTime    :", ti.TotalTime)
		log.Println("IsConnReused :", ti.IsConnReused)
		log.Println("IsConnWasIdle:", ti.IsConnWasIdle)
		log.Println("ConnIdleTime :", ti.ConnIdleTime)

		return nil
	})

	restyClient.
		SetHostURL(c.ApiHost).
		SetHeader("x-api-key", c.ApiKey).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "terraform-provider-paperspace").
		SetHeader("ps_client_name", "terraform-provider-paperspace")

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
	log.Printf("Response Received At: %v", resp.ReceivedAt())
	log.Printf("Response Body: %v", resp) // or resp.String() or string(resp.Body())
}
