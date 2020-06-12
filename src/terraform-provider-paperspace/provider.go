package main

import (
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFunc("PAPERSPACE_API_KEY"),
			},
			"api_host": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncAllowMissingDefault("PAPERSPACE_API_HOST", "https://api.paperspace.io"),
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: envDefaultFuncAllowMissing("PAPERSPACE_REGION"),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"paperspace_machine": resourceMachine(),
			"paperspace_script":  resourceScript(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"paperspace_network":  dataSourceNetwork(),
			"paperspace_template": dataSourceTemplate(),
			"paperspace_user":     dataSourceUser(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return nil, nil
	}
}

func envDefaultFuncAllowMissing(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		v := os.Getenv(k)
		return v, nil
	}
}

func envDefaultFuncAllowMissingDefault(k string, d string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		return d, nil
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := ClientConfig{
		APIKey:  d.Get("api_key").(string),
		APIHost: d.Get("api_host").(string),
		Region:  d.Get("region").(string),
	}

	log.Printf("[INFO] paperspace provider api_key %v", config.APIKey)
	log.Printf("[INFO] paperspace provider api_host %v", config.APIHost)
	if config.Region != "" {
		log.Printf("[INFO] paperspace provider region %v", config.Region)
	}

	return config.Client()
}
