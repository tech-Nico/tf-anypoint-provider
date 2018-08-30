package sdk

import (
	"github.com/tech-nico/terraform-provider-anypoint/anypoint"
	"golang.org/x/text/unicode/cldr"
)

type AccessManagement struct {
	AuthToken string

}

func (am *AccessManagement) NewAccessManagement(c anypoint.Config) (*AccessManagement){
	payload := make(map[string]string)
	payload["username"] = c.Username
	payload["password"] = c.Password


	c.HTTPClient.Post(c.Hostname, "application/json", payload)
}

func (am *AccessManagement) getAuthToken(username, password string) {

}