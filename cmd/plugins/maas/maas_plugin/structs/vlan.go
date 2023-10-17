package structs

// VlanStruct is the struct for a VLAN
type VlanStruct struct {
	Vid           int         `json:"vid" yaml:"vid"`
	Mtu           int         `json:"mtu" yaml:"mtu"`
	DhcpOn        bool        `json:"dhcp_on" yaml:"dhcp_on"`
	ExternalDhcp  interface{} `json:"external_dhcp" yaml:"external_dhcp"`
	RelayVlan     interface{} `json:"relay_vlan" yaml:"relay_vlan"`
	Id            int         `json:"id" yaml:"id"`
	PrimaryRack   string      `json:"primary_rack" yaml:"primary_rack"`
	Space         string      `json:"space" yaml:"space"`
	Fabric        string      `json:"fabric" yaml:"fabric"`
	FabricId      int         `json:"fabric_id" yaml:"fabric_id"`
	SecondaryRack string      `json:"secondary_rack" yaml:"secondary_rack"`
	Name          string      `json:"name" yaml:"name"`
	ResourceUri   string      `json:"resource_uri" yaml:"resource_uri"`
}
