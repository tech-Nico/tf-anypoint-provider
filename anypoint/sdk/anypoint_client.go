package sdk

import (
	"fmt"
)

type AnypointClient struct {
	Auth *Auth
}

func NewAnypointClient(uri string, username, password string, insecure bool) (*AnypointClient, error) {
	ac := new(AnypointClient)
	var err error
	ac.Auth, err = NewAuthWithCredentials(uri, username, password, insecure)

	if err != nil {
		return nil, fmt.Errorf("Error while creating a new instance of AnypointClient: %s", err)
	}

	return ac, nil
}
