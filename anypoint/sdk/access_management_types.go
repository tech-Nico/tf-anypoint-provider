package sdk

type AccessManagement struct {
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

type Entitlements struct {
	CreateSubOrgs      bool              `json:"createSubOrgs,omitempty"`
	GlobalDeployment   bool              `json:"globalDeployment,omitempty"`
	CreateEnvironments bool              `json:"createEnvironments,omitempty"`
	ProductionVCores   EntitlementStatus `json:"vCoresProduction,omitempty"`
	SandboxVCores      EntitlementStatus `json:"vCoresSandbox,omitempty"`
	DesignVCores       EntitlementStatus `json:"vCoresDesign,omitempty"`
	StaticIPs          EntitlementStatus `json:"staticIps,omitempty"`
	VPCs               EntitlementStatus `json:"vpcs,omitempty"`
	LoadBalancer       EntitlementStatus `json:"loadBalancer,omitempty"`
	VPNs               EntitlementStatus `json:"vpns,omitempty"`
}

type EntitlementStatus struct {
	Assigned float64 `json:"assigned"`
}
type BusinessGroup struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name,omitempty"`
	OwnerId      string       `json:"ownerId,omitempty"`
	ParentOrgId  string       `json:"parentOrganizationId,omitempty"`
	Entitlements Entitlements `json:"entitlements,omitempty"`
	ClientID     string       `json:"clientId:omitempty"`

	Domain                string          `json:"domain,omitempty"`
	ProviderID            string          `json:"idprovider_id,omitempty"`
	IsFederated           bool            `json:"isFederated,omitempty"`
	IsMaster              bool            `json:"isMaster,omitempty"`
	OwnerName             string          `json:"ownerName,omitempty"`
	ParentOrganizationIDs []string        `json:"parentOrganizationIds,omitempty:"`
	SessionTimeout        int             `json:"sessionTimeout,omitempty"`
	SubOrganizations      []BusinessGroup `json:"subOrganizations,omitempty"`
	TenantOrgIDs          []string        `json:"tenantOrganizationIds,omitempty"`
}

type Users struct {
	Total int    `json:"total,omitempty"`
	Data  []User `json:"data,omitempty"`
}

type User struct {
	ID                 string `json:"id,omitempty"`
	Username           string `json:"username,omitempty"`
	Firstname          string `json:"firstName,omitempty"`
	Lastname           string `json:"lastName,omitempty"`
	Email              string `json:"email,omitempty"`
	OrganizationID     string `json:"organizationId,omitempty"`
	Enabled            bool   `json:"enabled,omitempty"`
	IdentityProviderID string `json:"idprovider_id,omitempty"`
	Type               string `json:"type,omitempty"`
}
