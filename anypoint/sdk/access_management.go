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
	"fmt"
	"log"
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
	log.Printf("Logging in with %s", body)

	authToken := new(AuthToken)

	err := httpClient.POST(&body, LOGIN, &authToken)

	if err != nil {
		return "", fmt.Errorf("Error during login with user %q : %s", pUsername, err)
	}

	log.Printf("Been able to login. Auth Token:  %v", (*authToken).BearerToken)
	return authToken.BearerToken, nil
}

func (auth *Auth) Me() interface{} {
	log.Printf("Call to %s", ME)
	var res interface{}
	err := auth.client.GET(ME, &res)

	if err != nil {
		log.Fatalf("Error while retrieving user details: %s", err)
	}

	return res
}

func (auth *Auth) Hierarchy() BusinessGroup {
	me := auth.Me()
	data := me.(map[string]interface{})

	orgId := data["user"].(map[string]interface{})["organization"].(map[string]interface{})["id"].(string)
	path := hierarchyPath(orgId)

	var res = BusinessGroup{}
	err := auth.client.GET(path, &res)

	if err != nil {
		log.Fatalf("HTTP error while retrieving user details: %s", err)
	}

	return res
}



func (auth *Auth) GetBusinessGroupHierarchy(bgID string) (BusinessGroup, error) {
	path := hierarchyPath(bgID)
	var res BusinessGroup

	err := auth.client.GET(path, &res)

	if err != nil {
		return BusinessGroup{}, fmt.Errorf("HTTP error while retrieving user details: %s", err)
	}

	return res, nil
}

func (auth *Auth) GetBusinessGroupByID(bgID string) (BusinessGroup, error) {
	path := organizationPath(bgID)
	var res BusinessGroup

	err := auth.client.GET(path, &res)

	if err != nil {
		return BusinessGroup{}, fmt.Errorf("HTTP error while retrieving user details: %s", err)
	}

	return res, nil
}

func (auth *Auth) GetBusinessGroup(parentID, bgName string) (BusinessGroup, error) {

	bgStructure, err := auth.GetBusinessGroupHierarchy(parentID)

	if err != nil {
		return BusinessGroup{}, fmt.Errorf("error while retrieving hierarchi for business group with ID %s. Error : %s", parentID, err)
	}

	if bgStructure.SubOrganizations == nil || len(bgStructure.SubOrganizations) <= 0 {
		return BusinessGroup{}, nil
	}

	for _, subOrg := range bgStructure.SubOrganizations {
		if subOrg.Name == bgName {
			return subOrg, nil
		}
	}

	return BusinessGroup{}, nil

}


// FindBusinessGroup search for the given business group (specified in the format "Parent\Child\Grand-Nephew") and
//return its ID

func (auth *Auth) FindBusinessGroup(path string) (string, error) {
	currentOrgId := ""

	groups := auth.CreateBusinessGroupPath(path)

	hierarchy := auth.Hierarchy()

	subOrganizations := hierarchy.SubOrganizations

	if len(groups) == 1 {
		return hierarchy.ID, nil
	}

	for _, currGroup := range groups {
		for organization := 0; organization < len(subOrganizations); organization++ {
			currOrg := subOrganizations[organization]

			if currOrg.Name == currGroup {
				currentOrgId = currOrg.ID
				log.Printf("The matched org name is: %s", currOrg.Name)
				subOrganizations = currOrg.SubOrganizations
			}
		}
	}

	if currentOrgId == "" {
		return "", fmt.Errorf("cannot find business group %s", path)
	}

	return currentOrgId, nil
}

func (auth *Auth) CreateBusinessGroupPath(businessGroup string) []string {
	if businessGroup == "" {
		return make([]string, 0)
	}

	groups := []string{}
	group := ""
	pos := 0
	for ; pos < len(businessGroup)-1; pos++ {
		currChar := businessGroup[pos]
		if currChar == '/' {
			// Double backslash maps to business group with one backslash
			if businessGroup[pos+1] == '/' {
				group += "/"
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

func (auth *Auth) FindUserByUsername(orgId, username string) (*User, error) {

	params := make(map[string]string)
	params["limit"] = "20"
	params["offset"] = "0"
	params["search"] = username

	var response Users
	err := auth.client.GETWithParams(searchUserPath(orgId), params, &response)

	if err != nil {
		return nil, fmt.Errorf("error while searching for user with username %s : %s", username, err)
	}
	if response.Total > 1 {
		return nil, fmt.Errorf("%d results returned while searching for user %s", response.Total, username)
	}

	return &response.Data[0], nil
}

func (auth *Auth) UpdateBusinessGroup() {

}

//Create a new business group under the given business group ID. Returns the newly created business group ID
func (auth *Auth) CreateBusinessGroup(ownerUsername, parentBGID, newBGName string, entitlements Entitlements) (BusinessGroup, error) {

	user, err := auth.FindUserByUsername(parentBGID, ownerUsername)

	if err != nil {
		return BusinessGroup{}, fmt.Errorf("error when searching for owner [%s] of new bg [%s] to be created: %s", ownerUsername, newBGName, err)
	}

	//Find the business group to check it doesn't exist already. If so we need to update instead.
	newBG, err := auth.GetBusinessGroup(parentBGID, newBGName)

	if err != nil {
		return BusinessGroup{}, fmt.Errorf("error while creating business group %s : %s", newBGName, err)
	}

	var response BusinessGroup

	if newBG.ID == "" {

		newBG = BusinessGroup{
			Name:         newBGName,
			OwnerId:      user.ID,
			Entitlements: entitlements,
			ParentOrgId:  parentBGID,
		}
		log.Printf("Creating new business group [%s]", newBGName)

		err = auth.client.POST(newBG, "/accounts/api/organizations", &response)

	} else {
		newBG.Name = newBGName
		newBG.OwnerId = user.ID
		newBG.Entitlements = entitlements
		newBG.ParentOrgId = parentBGID

		log.Printf("Updating BG %s [%s]", newBG.Name, newBG.ID)
		err = auth.client.PUT(newBG, "/accounts/api/organizations/"+newBG.ID, &response)
	}

	if err != nil {
		return BusinessGroup{}, fmt.Errorf("Error while creating/updating business group %s : %s", newBGName, err)
	}

	return newBG, nil
}
