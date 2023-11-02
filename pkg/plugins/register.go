package plugins

// RegisterArgs is a struct to pass in any configuration parameters from the plugin host
type RegisterArgs struct {
	DefaultSecure bool
	Config        map[string]interface{}
}

// RegisterReply is a struct to pass back any endpoints that the plugin wants to register
type RegisterReply struct {
	Endpoints []*PluginEndpoint
}
