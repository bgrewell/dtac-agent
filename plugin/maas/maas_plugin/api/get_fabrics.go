package api

import (
	"encoding/json"
	"github.com/BGrewell/dtac-agent/plugin/maas/maas_plugin/structs"
)

func GetFabrics(settings *structs.MAASSettings) ([]*structs.Fabric, error) {
	endpoint := "fabrics/"

	body, err := Get(endpoint, settings)
	if err != nil {
		return nil, err
	}

	var fabrics []*structs.Fabric
	err = json.Unmarshal(body, &fabrics)
	if err != nil {
		return nil, err
	}

	return fabrics, nil
}
