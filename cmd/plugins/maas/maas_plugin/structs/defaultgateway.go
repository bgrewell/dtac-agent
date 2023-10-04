package structs

type DefaultGatewayStruct struct {
	Ipv4 map[string]interface{} `json:"ipv4" yaml:"ipv4"`
	Ipv6 map[string]interface{} `json:"ipv6" yaml:"ipv6"`
}
