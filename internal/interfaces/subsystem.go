package interfaces

import "github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"

// Subsystem is the interface for the subsystems
type Subsystem interface {
	Endpoints() []endpoint.Endpoint
	Enabled() bool
	Name() string
}
