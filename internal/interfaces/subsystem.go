package interfaces

// Subsystem is the interface for the subsystems
type Subsystem interface {
	Register() error
	Enabled() bool
	Name() string
}
