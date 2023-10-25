package interfaces

import "context"

// APIAdapter is the interface for the frontend APIs
type APIAdapter interface {
	Register(subsystems []Subsystem) (err error)
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
	Name() string
}
