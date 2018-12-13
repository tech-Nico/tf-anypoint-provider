package sdk

import (
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	username, password := getCredentials(t)
	auth, err := NewAuthWithCredentials("https://anypoint.mulesoft.com", username, password, true, true)

	if err != nil {
		t.Errorf("Error while logging in into Anypoint: %s", err)
	}

	t.Logf("Obtained auth object: %v\n", auth)

	if auth.GetAuthenticatedHttpClient() == nil {
		t.Error("Error while testing instatiating a new Auth object. Rest client is nil")
	}

}

func TestAuth_FindBusinessGroup(t *testing.T) {
	bgPath := "RootOrg\\Sub Org 1"
	auth := getAuth(t)
	id, err := auth.FindBusinessGroup(bgPath)

	if err != nil {
		t.Fatalf("Error while searching for business group [%s] : %s", bgPath, err)
	}

	t.Logf("Got expected result when searching for business group [%s] : [id=%s]", bgPath, id)

}

func TestAuth_CreateBusinessGroup_StepByStep(t *testing.T) {
	username, _ := getCredentials(t)
	auth := getAuth(t)
	parentBgPath := "RootOrg"
	newBGName := "TestAuth_CreateBusinessGroup_StepByStep"
	parentBgID, err := auth.FindBusinessGroup(parentBgPath)
	if err != nil {
		t.Fatalf("Unable to find parent org [%s] when testing whether I can create a new BG %s : %s", parentBgPath, newBGName, err)
	}
	ents := Entitlements{
		CreateEnvironments: true,
	}
	newBG, err := auth.CreateBusinessGroup(username, parentBgID, newBGName, ents)

	if err != nil {
		t.Fatalf("Error while creating new BG [%s\\%s] -> %s", parentBgPath, newBGName, err)
	}

	t.Logf("Successful in creating BG [%s\\%s]: new BG id is [%s]", parentBgPath, newBGName, newBG.ID)
}

func TestAuth_CreateBusinessGroup(t *testing.T) {
	username, password := getCredentials(t)
	anypoint, err := NewAnypointClient("https://anypoint.mulesoft.com", username, password, true, true)

	if err != nil {
		t.Fatalf("Error creating a new instance of AnypointClient: %s", err)
	}

	parentBgPath := "RootOrg"
	newBGName := "TestAuth_CreateBusinessGroup"
	parentBgID, err := anypoint.Auth.FindBusinessGroup(parentBgPath)

	if err != nil {
		t.Fatalf("Unable to find parent org [%s] when testing whether I can create a new BG %s : %s", parentBgPath, newBGName, err)
	}
	ents := Entitlements{
		CreateEnvironments: true,
	}
	newBG, err := anypoint.Auth.CreateBusinessGroup(username, parentBgID, newBGName, ents)

	if err != nil {
		t.Fatalf("Error while creating new BG [%s\\%s] -> %s", parentBgPath, newBGName, err)
	}

	t.Logf("Successful in creating BG [%s\\%s]: new BG id is [%s]", parentBgPath, newBGName, newBG.ID)
}

func TestAuth_FindUserByUsername(t *testing.T) {
	username, _ := getCredentials(t)
	auth := getAuth(t)
	orgId, err := auth.FindBusinessGroup("RootOrg")
	if err != nil {
		t.Fatalf("Error while searching for business group RootOrg: %s", err)
	}

	user, err := auth.FindUserByUsername(orgId, username)
	if err != nil {
		t.Fatalf("Error while searching for user %s in org RootOrg: %s", username, err)
	}

	t.Logf("Test successful. User %s found: %s %s - %s", username, user.Firstname, user.Lastname, user.Email)
}
func getCredentials(t *testing.T) (string, string) {
	username := os.Getenv("ANYPOINT_USERNAME")
	password := os.Getenv("ANYPOINT_PASSWORD")

	if username == "" {
		t.Fatalf("ANYPOINT_USERNAME environment variable not set")
	}

	if password == "" {
		t.Fatalf("ANYPOINT_PASSWORD environment variable not set")
	}

	return username, password
}

func getAuth(t *testing.T) *Auth {
	username, password := getCredentials(t)
	auth, err := NewAuthWithCredentials("https://anypoint.mulesoft.com", username, password, true, true)

	if err != nil {
		t.Errorf("Error while logging in into Anypoint: %s", err)
	}

	return auth

}
