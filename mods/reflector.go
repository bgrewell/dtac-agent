package mods

type Reflector interface {
	Proto() string
	Port() int
	SetPort(int)
	Start() error
	Stop() error
}
