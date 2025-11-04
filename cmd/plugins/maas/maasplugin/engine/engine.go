package engine

import (
	"errors"

	"github.com/bgrewell/dtac-agent/cmd/plugins/maas/maasplugin/api"
	"github.com/bgrewell/dtac-agent/cmd/plugins/maas/maasplugin/structs"
)

// Engine is the main engine for the MAAS plugin
type Engine struct {
	Settings *structs.MAASSettings
	running  bool
	errored  bool
	err      error
	machines []*structs.Machine
	fabrics  []*structs.Fabric
}

func (e *Engine) CreateMachine() (results []byte, err error) {
	// TODO: Need to think of where the dividing lines are in this package, for example the engine probably shouldn't
	//  be taking in grpc structs and working with them, but rather the engine should be working with the structs
	//  defined in the maasplugin package. The engine should be agnostic to the grpc structs.
	return []byte{}, errors.New("Not implemented")
}

// Machines returns a list of machines from the MAAS server
func (e *Engine) Machines() []*structs.Machine {
	return e.machines
}

// Fabrics returns a list of fabrics from the MAAS server
func (e *Engine) Fabrics() []*structs.Fabric {
	return e.fabrics
}

// Running returns true if the engine is running
func (e *Engine) Running() bool {
	return e.running
}

// Failed returns true if the engine has errored
func (e *Engine) Failed() bool {
	return e.errored
}

// ErrDetails returns the error details
func (e *Engine) ErrDetails() error {
	return e.err
}

// Start starts the engine
func (e *Engine) Start() error {
	e.running = true
	go func() {
		var err error
		for e.running {

			// Update machines
			e.machines, err = api.GetMachines(e.Settings)
			if err != nil {
				e.running = false
				e.errored = true
				e.err = err
			}

			// Update fabrics
			e.fabrics, err = api.GetFabrics(e.Settings)
			if err != nil {
				e.running = false
				e.errored = true
				e.err = err
			}
		}
	}()
	return nil
}

// Stop stops the engine
func (e *Engine) Stop() error {
	e.running = false
	return nil
}
