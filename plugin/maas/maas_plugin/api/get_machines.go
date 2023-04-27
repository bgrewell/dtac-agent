package api

import (
	"encoding/json"
	"github.com/BGrewell/dtac-agent/plugin/maas/maas_plugin/structs"
)

func GetMachines(settings *structs.MAASSettings) ([]*structs.Machine, error) {
	endpoint := "machines/"

	body, err := Get(endpoint, settings)
	if err != nil {
		return nil, err
	}

	var machines []*structs.Machine
	err = json.Unmarshal(body, &machines)
	if err != nil {
		return nil, err
	}

	return machines, nil
}
