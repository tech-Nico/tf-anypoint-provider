package anypoint

import (
	"github.com/hashicorp/go-cleanhttp"
	"net/http"
	"crypto/tls"
	"github.com/tech-nico/terraform-provider-anypoint/anypoint/sdk"
)

type Config struct {
	Username string
	Password string
	Hostname string
	Insecure bool
	HTTPClient *http.Client
}

func (c *Config) Authenticate() (error){
	c.HTTPClient = cleanhttp.DefaultClient()

	if c.Insecure {
		transport := c.HTTPClient.Transport.(*http.Transport)
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	am := sdk.NewAccessManagement(*c)

}