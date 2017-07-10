package main

import (
  "github.com/hashicorp/terraform/helper/schema"
  "log"
  "os"
)

func Provider() *schema.Provider {
  return &schema.Provider{
    Schema: map[string]*schema.Schema{
      "apiKey": &schema.Schema{
        Type:        schema.TypeString,
        Optional:    true,
        DefaultFunc: envDefaultFunc("PAPERSPACE_API_KEY"),
      },
      "apiHost": &schema.Schema{
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
      "paperspace_machine":  resourceMachine(),
      "paperspace_script":   resourceScript(),
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
	config := Config{
		ApiKey: d.Get("apiKey").(string),
		ApiHost: d.Get("apiHost").(string),
		Region: d.Get("region").(string),
	}

  log.Printf("[INFO] paperspace provider apiKey %v", config.ApiKey)
  log.Printf("[INFO] paperspace provider apiHost %v", config.ApiHost)
  if config.Region != "" {
    log.Printf("[INFO] paperspace provider region %v", config.Region)
  }

	return config.Client()
}
