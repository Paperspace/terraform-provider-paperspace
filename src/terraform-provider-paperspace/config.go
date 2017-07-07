package main

import (
  "gopkg.in/resty.v0"
  "log"
)

type Config struct {
	ApiKey  string
	ApiHost string
	Region  string
}

type PaperspaceClient struct {
  ApiKey  string
	ApiHost string
	Region  string
	RestyClient *resty.Client
}

func (c *Config) Client() (PaperspaceClient, error) {

  restyClient := resty.New();

  restyClient.
  SetHostURL(c.ApiHost).
  SetHeader("x-api-key", c.ApiKey).
  SetHeader("Accept", "application/json")

	client := PaperspaceClient{
		ApiKey:  c.ApiKey,
		ApiHost: c.ApiHost,
		Region:  c.Region,
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
  log.Printf("Response Body: %v", resp) // or resp.String() or string(resp.Body())
}
