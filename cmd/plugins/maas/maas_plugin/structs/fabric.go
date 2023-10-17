package structs

// Fabric is the struct for a fabric
type Fabric struct {
	ClassType   interface{}  `json:"class_type" yaml:"class_type"`
	Vlans       []VlanStruct `json:"vlans" yaml:"vlans"`
	Id          int          `json:"id" yaml:"id"`
	Name        string       `json:"name" yaml:"name"`
	ResourceUri string       `json:"resource_uri" yaml:"resource_uri"`
}
