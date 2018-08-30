// Copyright © 2017 Nico Balestra <functions@protonmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Auth struct {
	uri      string
	insecure bool
	client   *RestClient
	Token    string
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthToken struct {
	BearerToken string `json:"access_token,omitempty"`
}

const (
	LOGIN     string = "/accounts/login"
	ME        string = "/accounts/api/me"
	HIERARCHY string = "/accounts/api/organizations/{orgId}/hierarchy"
)

func NewAuthWithCredentials(uri, username, password string, insecure bool) (*Auth, error) {
	client := NewRestClient(uri, insecure)
	token, err := login(client, username, password)
	if err != nil {
		return nil, fmt.Errorf("Error while logging in into Anypoint Platform: %s", err)
	}

	client.AddAuthHeader(token)
	return &Auth{
		uri,
		insecure,
		client,
		token,
	}, nil
}

func (auth *Auth) GetAuthenticatedHttpClient() *RestClient {
	return auth.client
}

/**
Return an HTTP Client which will inject the necessary headers to authenticate calls against ARM
*/
func (auth *Auth) GetARMAuthenticatedHttpClient(orgId, envId string) *RestClient {
	//We are not caching here since we could call this function with different orgId and envId in the same execution
	armClient := NewRestClient(auth.uri, auth.insecure)
	armClient = NewRestClient(auth.uri, auth.insecure)
	armClient.AddAuthHeader(auth.Token)
	armClient.AddEnvHeader(envId)
	armClient.AddOrgHeader(orgId)
	return armClient
}

// Login the given user and return the bearer Token
func login(httpClient *RestClient, pUsername, pPassword string) (string, error) {
	body := LoginPayload{
		Username: pUsername,
		Password: pPassword,
	}

	authToken := new(AuthToken)

	_, err := httpClient.POST(body, LOGIN, Application_Json, &authToken)

	if err != nil {
		return "", fmt.Errorf("Error during login with user %q : %s", pUsername, err)
	}

	log.Print("Been able to login: ", *authToken)
	return authToken.BearerToken, nil
}

func (auth *Auth) Me() []byte {
	log.Printf("Call to %s", ME)
	resp, err := auth.client.GET(ME)

	if err != nil {
		log.Fatalf("Error while retrieving user details: %s", err)
	}

	return resp
}

func (auth *Auth) Hierarchy() []byte {
	me := auth.Me()
	var data map[string]interface{}
	if err := json.Unmarshal(me, &data); err != nil {
		log.Fatalf("Invalid JSON response when retrieving user details: %s", err)
	}

	orgId := data["user"].(map[string]interface{})["organization"].(map[string]interface{})["id"].(string)
	path := strings.Replace(HIERARCHY, "{orgId}", orgId, -1)

	resp, err := auth.client.GET(path)

	if err != nil {
		log.Fatalf("HTTP error while retrieving user details: %s", err)
	}

	return resp
}

func (auth *Auth) FindBusinessGroup(path string) string {
	currentOrgId := ""

	groups := auth.createBusinessGroupPath(path)

	var data map[string]interface{}
	hierarchy := auth.Hierarchy()

	if err := json.Unmarshal(hierarchy, &data); err != nil {
		log.Fatalf("Error while querying for hierarchy : %s", err)
	}

	subOrganizations := data["subOrganizations"].([]interface{})

	if len(groups) == 1 {
		return data["id"].(string)
	}

	for _, currGroup := range groups {
		for organization := 0; organization < len(subOrganizations); organization++ {
			jsonObject := subOrganizations[organization].(map[string]interface{})

			if jsonObject["name"].(string) == currGroup {
				currentOrgId = jsonObject["id"].(string)
				log.Printf("The matched org name is: %s", jsonObject["name"].(string))
				subOrganizations = jsonObject["subOrganizations"].([]interface{})
			}
		}
	}

	if currentOrgId == "" {
		log.Fatalf("Cannot find business group %s", path)
	}

	return currentOrgId
}

func (auth *Auth) createBusinessGroupPath(businessGroup string) []string {
	if businessGroup == "" {
		return make([]string, 0)
	}

	groups := []string{}
	group := ""
	pos := 0
	for ; pos < len(businessGroup)-1; pos++ {
		currChar := businessGroup[pos]
		if currChar == '\\' {
			// Double backslash maps to business group with one backslash
			if businessGroup[pos+1] == '\\' {
				group += "\\"
				pos++
				// Single backslash starts a new business group
			} else {
				groups = append(groups, group)
				group = ""
			}
			// Non backslash characters are mapped to the group
		} else {
			group += string(currChar)
		}
	}

	if pos < len(businessGroup) { // Do not end with backslash {
		group += string(businessGroup[len(businessGroup)-1])
	}
	groups = append(groups, string(group))

	return groups
}
