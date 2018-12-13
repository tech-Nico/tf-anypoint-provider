package sdk

import (
	"fmt"
)

type AnypointClient struct {
	AccessManagement *AccessManagement
}

func NewAnypointClient(uri string, username, password string, insecure, httpWireLog bool) (*AnypointClient, error) {
	ac := new(AnypointClient)
	var err error
	ac.AccessManagement, err = NewAuthWithCredentials(uri, username, password, insecure, httpWireLog)

	if err != nil {
		return nil, fmt.Errorf("Error while creating a new instance of AnypointClient: %s", err)
	}

	return ac, nil
}
