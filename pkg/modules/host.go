package modules

import (
	"github.com/bgrewell/dtac-agent/pkg/modules/utility"
)

// ModuleHost is the interface for the module host
type ModuleHost interface {
	Serve() error
	GetPort() int
}

// NewModuleHost creates a new ModuleHost
func NewModuleHost(module Module) (hostModule ModuleHost, err error) {
	key := utility.NewRandomSymmetricKey()
	mod := &DefaultModuleHost{
		Module:     module,
		Proto:      "tcp",
		IP:         "127.0.0.1",
		APIVersion: "mod_api_1.0",
		encryptor:  utility.NewRPCEncryptor(key),
	}

	return mod, nil
}
