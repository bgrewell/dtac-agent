package network

type Route interface {
	String() string
	JSON() string
	Create() error
	Update() error
	Remove() error
	Applied() bool
}
