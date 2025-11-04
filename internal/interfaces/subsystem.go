package interfaces

import (
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// Subsystem is the interface for the subsystems
type Subsystem interface {
	Endpoints() []*endpoint.Endpoint
	Enabled() bool
	Name() string
}
