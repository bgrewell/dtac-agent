package engine

import (
	api2 "github.com/bgrewell/dtac-agent/cmd/plugins/maas/maasplugin/api"
	structs2 "github.com/bgrewell/dtac-agent/cmd/plugins/maas/maasplugin/structs"
)

// Engine is the main engine for the MAAS plugin
type Engine struct {
	Settings *structs2.MAASSettings
	running  bool
	errored  bool
	err      error
	machines []*structs2.Machine
	fabrics  []*structs2.Fabric
}

// Machines returns a list of machines from the MAAS server
func (e *Engine) Machines() []*structs2.Machine {
	return e.machines
}

// Fabrics returns a list of fabrics from the MAAS server
func (e *Engine) Fabrics() []*structs2.Fabric {
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
			e.machines, err = api2.GetMachines(e.Settings)
			if err != nil {
				e.running = false
				e.errored = true
				e.err = err
			}

			// Update fabrics
			e.fabrics, err = api2.GetFabrics(e.Settings)
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
