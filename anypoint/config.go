package anypoint

import (
	"fmt"
	"github.com/tech-nico/terraform-provider-anypoint/anypoint/sdk"
)

type Config struct {
	Username       string
	AnypointClient *sdk.AnypointClient
}

func NewConfig(hostname, username, password string, insecureSSL, httpWireLog bool) (*Config, error) {

	anypointClient, err := sdk.NewAnypointClient(hostname, username, password, insecureSSL, httpWireLog)

	if err != nil {
		return nil, fmt.Errorf("Error while creating an instance of AnypointClient : %s", err)
	}

	config := &Config{
		Username:       username,
		AnypointClient: anypointClient,
	}

	return config, nil
}
