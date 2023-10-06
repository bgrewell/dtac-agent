package interfaces

type Subsystem interface {
	Register() error
	Enabled() bool
	Name() string
}
