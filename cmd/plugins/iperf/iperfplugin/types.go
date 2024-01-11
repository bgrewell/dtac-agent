package iperfplugin

// IperfTestIDRequest is the struct for the validation of functions that require an id parameter
type IperfTestIDRequest struct {
	ID []string `json:"id"`
}

// IperfServerStartRequest is a struct that is used to validate server start requests
type IperfServerStartRequest struct {
	BindAddr []string `json:"bind_addr,omitempty"`
}

// IperfClientStartRequest is a struct that is used to validate client start requests
type IperfClientStartRequest struct {
	Host []string `json:"host"`
}
