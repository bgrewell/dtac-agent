package api

import (
	"encoding/json"
	structs2 "github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maasplugin/structs"
)

// GetMachines returns a list of machines from the MAAS server
func GetMachines(settings *structs2.MAASSettings) ([]*structs2.Machine, error) {
	endpoint := "machines/"

	body, err := Get(endpoint, settings)
	if err != nil {
		return nil, err
	}

	var machines []*structs2.Machine
	err = json.Unmarshal(body, &machines)
	if err != nil {
		return nil, err
	}

	return machines, nil
}
