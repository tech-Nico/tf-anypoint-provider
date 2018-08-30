package sdk

type Endpoint struct {
	Id                   int    `json:"id"`
	OrgID                string `json:"masterOrganizationId"`
	ApiID                int
	VersionID            int    `json:"apiVersionId"`
	Type                 string `json:"type"`
	Uri                  string `json:"uri"`
	ProxyUri             string `json:"proxyUri"`
	ProxyRegistrationUri string `json:"proxyRegistrationUri"`
	IsCloudHub           bool   `json:"isCloudHub"`
	ReferencesUserDomain bool   `json:"referencesUserDomain"`
	ResponseTimeout      int    `json:"responseTimeout"`
}

type ClusterServer struct {
	ID   float64 `yaml:"id,omitempty" json:"serverId"`
	Name string  `yaml:"name" json:"-"`
	Ip   string  `yaml:"ip" json:"serverIp,omitempty"`
}

type Cluster struct {
	ClusterName string          `yaml:"cluster_name" json:"name"`
	Multicast   bool            `yaml:"mulsticast" json:"multicastEnabled"`
	Servers     []ClusterServer `yaml:"servers" json:"servers"`
}

type Clusters struct {
	ACluster []Cluster `yaml:"clusters"`
}
