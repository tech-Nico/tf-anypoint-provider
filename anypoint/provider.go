package anypoint

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tech-nico/terraform-provider-anypoint/anypoint/sdk"
	"github.com/tech-nico/terraform-provider-anypoint/anypoint/sdk/rest"
	"net/url"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"hostname": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true, //Default to anypoint.mulesoft.com
				DefaultFunc:  schema.EnvDefaultFunc("ANYPOINT_URL", "anypoint.mulesoft.com"),
				ValidateFunc: validateUrl,
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_PASSWORD", nil),
			},
			"insecureSSL": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: false,
				Default:  false,
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"anypoint_business_group": resourceBusinessGroup(),
		},
	}
}

func validateUrl(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}

	urlStr := v.(string)
	_, err := url.Parse(urlStr)

	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not a valid URL: %s", urlStr, err))
	}

	return
}

func providerConfigure(rd *schema.ResourceData) (interface{}, error) {

	hostname := rd.Get("hostname").(string)
	username := rd.Get("username").(string)
	password := rd.Get("password").(string)
	insecure := rd.Get("insecure").(bool)

	anypointClient, err := sdk.NewAnypointClient(hostname, username, password, insecure)

	if err != nil {
		return nil, fmt.Errorf("Error while creating an instance of AnypointClient : %s", err)
	}

	config := Config{
		AnypointClient: anypointClient,
	}

	return &config, nil
}
