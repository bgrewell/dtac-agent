package api

import (
	"encoding/json"
	structs2 "github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maas_plugin/structs"
)

func GetFabrics(settings *structs2.MAASSettings) ([]*structs2.Fabric, error) {
	endpoint := "fabrics/"

	body, err := Get(endpoint, settings)
	if err != nil {
		return nil, err
	}

	var fabrics []*structs2.Fabric
	err = json.Unmarshal(body, &fabrics)
	if err != nil {
		return nil, err
	}

	return fabrics, nil
}
