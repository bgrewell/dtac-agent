package network

type ICapabilities interface {
	Name() string
	Capability() string
}
