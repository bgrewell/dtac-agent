package structs

// DefaultGatewayStruct is the struct for the default gateway
type DefaultGatewayStruct struct {
	Ipv4 map[string]interface{} `json:"ipv4" yaml:"ipv4"`
	Ipv6 map[string]interface{} `json:"ipv6" yaml:"ipv6"`
}
