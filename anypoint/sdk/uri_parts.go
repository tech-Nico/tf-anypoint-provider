package sdk

import "strings"

const (
	BASE_URI     = "/accounts/api"
	LOGIN        = "/accounts/login"
	ME           = BASE_URI + "/me"
	ORGANIZATION = BASE_URI + "/organizations/{orgId}"
	HIERARCHY    = ORGANIZATION + "/hierarchy"
	SEARCH_USER  = ORGANIZATION + "/members"
)

func hierarchyPath(orgId string) string {
	return strings.Replace(HIERARCHY, "{orgId}", orgId, -1)
}

func searchUserPath(orgId string) string {
	return strings.Replace(SEARCH_USER, "{orgId}", orgId, -1)
}


func organizationPath(orgId string) string {
	return strings.Replace(ORGANIZATION, "{orgId}", orgId, -1)
}
