package engine

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/plugin/maas/maas_plugin/api"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/plugin/maas/maas_plugin/structs"
)

type Engine struct {
	Settings *structs.MAASSettings
	running  bool
	errored  bool
	err      error
	machines []*structs.Machine
	fabrics  []*structs.Fabric
}

func (e *Engine) Machines() []*structs.Machine {
	return e.machines
}

func (e *Engine) Fabrics() []*structs.Fabric {
	return e.fabrics
}

func (e *Engine) Running() bool {
	return e.running
}

func (e *Engine) Failed() bool {
	return e.errored
}

func (e *Engine) ErrDetails() error {
	return e.err
}

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

func (e *Engine) Stop() error {
	e.running = false
	return nil
}
