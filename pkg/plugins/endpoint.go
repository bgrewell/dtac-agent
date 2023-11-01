package plugins

import "github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"

type PluginEndpoint struct {
	*endpoint.Endpoint
	FunctionName string
}
