package anypoint

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"anypoint_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true, //Default to anypoint.mulesoft.com
				DefaultFunc:  schema.EnvDefaultFunc("ANYPOINT_URL", "https://anypoint.mulesoft.com"),
				ValidateFunc: validateUrl,
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Sensitive:   true,
				Optional:    false,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ANYPOINT_PASSWORD", nil),
			},
			"insecure_ssl": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"http_debug_log": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"ap_bg": resourceBusinessGroup(),
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

	hostname := rd.Get("anypoint_url").(string)
	username := rd.Get("username").(string)
	password := rd.Get("password").(string)
	insecure := rd.Get("insecure_ssl").(bool)
	httpWire := rd.Get("http_debug_log").(bool)

	return NewConfig(hostname, username, password, insecure, httpWire)

}
