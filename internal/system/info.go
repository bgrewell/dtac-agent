package system

import "go.uber.org/zap"

type SystemInfo struct {
	ProductName            string `json:"product_name"`
	OperatingSystemName    string `json:"operating_system_name"`
	OperatingSystemVersion string `json:"operating_system_version"`
}

func (si *SystemInfo) Initialize(log *zap.Logger) {
	pn, err := GetSystemProductName()
	if err != nil {
		log.Error("failed to get product name", zap.Error(err))
		si.ProductName = "unknown"
	} else {
		si.ProductName = pn
	}

	os, err := GetOSName()
	if err != nil {
		log.Error("failed to get os name", zap.Error(err))
		si.OperatingSystemName = "unknown"
	} else {
		si.OperatingSystemName = os
	}

	ver, err := GetOSVersion()
	if err != nil {
		log.Error("failed to get os version", zap.Error(err))
		si.OperatingSystemVersion = "unknown"
	} else {
		si.OperatingSystemVersion = ver
	}
}
