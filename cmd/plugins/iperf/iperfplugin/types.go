package iperfplugin

// IperfTestIDRequest is the struct for the validation of functions that require an id parameter
type IperfTestIDRequest struct {
	Id []string `json:"id"`
}

type IperfServerStartRequest struct {
	BindAddr []string `json:"bind_addr,omitempty"`
}

type IperfClientStartRequest struct {
	Host []string `json:"host"`
}
