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
		log.Println("[INFO] Request Info:")
		log.Println("[INFO] client.HostURL", c.HostURL)
		log.Println("[INFO] client.Header", c.Header)
		log.Println("[INFO] req.Method", req.Method)
		log.Println("[INFO] req.URL", req.URL)
		log.Println("[INFO] req.Body", req.Body)
		log.Println("[INFO] req.AuthScheme", req.AuthScheme)
		log.Println("[INFO] req.RawRequest", req.RawRequest)
		log.Println("[INFO] req.Error", req.Error)

		return nil
	})

	restyClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {

		// Explore response object
		log.Println("[INFO] Response Info:")
		// log.Println("[INFO] Error      :", err)
		log.Println("[INFO] Status Code:", resp.StatusCode())
		log.Println("[INFO] Status     :", resp.Status())
		log.Println("[INFO] Proto      :", resp.Proto())
		log.Println("[INFO] Time       :", resp.Time())
		log.Println("[INFO] Received At:", resp.ReceivedAt())
		log.Println("[INFO] Body       :\n", resp)

		// Explore trace info
		log.Println("[INFO] Request Trace Info:")
		ti := resp.Request.TraceInfo()
		log.Println("[INFO] DNSLookup    :", ti.DNSLookup)
		log.Println("[INFO] ConnTime     :", ti.ConnTime)
		log.Println("[INFO] TCPConnTime  :", ti.TCPConnTime)
		log.Println("[INFO] TLSHandshake :", ti.TLSHandshake)
		log.Println("[INFO] ServerTime   :", ti.ServerTime)
		log.Println("[INFO] ResponseTime :", ti.ResponseTime)
		log.Println("[INFO] TotalTime    :", ti.TotalTime)
		log.Println("[INFO] IsConnReused :", ti.IsConnReused)
		log.Println("[INFO] IsConnWasIdle:", ti.IsConnWasIdle)
		log.Println("[INFO] ConnIdleTime :", ti.ConnIdleTime)

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
	log.Printf("[INFO] Request: %v", reqDesc)
	log.Printf("[INFO] Error: %v", err)
	log.Printf("[INFO] Response Status Code: %v", resp.StatusCode())
	log.Printf("[INFO] Response Status: %v", resp.Status())
	log.Printf("[INFO] Response Time: %v", resp.Time())
	log.Printf("[INFO] Response Received At: %v", resp.ReceivedAt())
	log.Printf("[INFO] Response Body: %v", resp) // or resp.String() or string(resp.Body())
}
