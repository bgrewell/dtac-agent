package structs

type InterfaceStruct struct {
	Params          interface{}              `json:"params" yaml:"params"`
	Vendor          string                   `json:"vendor" yaml:"vendor"`
	Parents         []map[string]interface{} `json:"parents" yaml:"parents"`
	Product         string                   `json:"product" yaml:"product"`
	Type            interface{}              `json:"type" yaml:"type"`
	SystemId        string                   `json:"system_id" yaml:"system_id"`
	Tags            []string                 `json:"tags" yaml:"tags"`
	Links           []map[string]interface{} `json:"links" yaml:"links"`
	SriovMaxVf      int                      `json:"sriov_max_vf" yaml:"sriov_max_vf"`
	EffectiveMtu    int                      `json:"effective_mtu" yaml:"effective_mtu"`
	Name            string                   `json:"name" yaml:"name"`
	NumaNode        int                      `json:"numa_node" yaml:"numa_node"`
	MacAddress      string                   `json:"mac_address" yaml:"mac_address"`
	Id              int                      `json:"id" yaml:"id"`
	LinkConnected   bool                     `json:"link_connected" yaml:"link_connected"`
	Enabled         bool                     `json:"enabled" yaml:"enabled"`
	Discovered      interface{}              `json:"discovered" yaml:"discovered"`
	Children        []map[string]interface{} `json:"children" yaml:"children"`
	LinkSpeed       int                      `json:"link_speed" yaml:"link_speed"`
	InterfaceSpeed  int                      `json:"interface_speed" yaml:"interface_speed"`
	Vlan            map[string]interface{}   `json:"vlan" yaml:"vlan"`
	FirmwareVersion string                   `json:"firmware_version" yaml:"firmware_version"`
	ResourceUrl     string                   `json:"resource_url" yaml:"resource_url"`
}
