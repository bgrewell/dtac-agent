package system

import (
	"go.uber.org/zap"
)

// Info is the struct for the system information
type Info struct {
	UUID                   string `json:"uuid"`
	ProductName            string `json:"product_name"`
	OperatingSystemName    string `json:"operating_system_name"`
	OperatingSystemVersion string `json:"operating_system_version"`
}

// Initialize initializes the system info
func (si *Info) Initialize(log *zap.Logger) {
	pn, err := GetSystemProductName()
	if err != nil {
		log.Error("failed to get product name", zap.Error(err))
		si.ProductName = "unknown"
	} else {
		si.ProductName = pn
	}

	id, err := GetSystemUUID()
	if err != nil {
		log.Error("failed to get system uuid", zap.Error(err))
		si.UUID = "unknown"
	} else {
		si.UUID = id
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

func (si *Info) serializeOs() (os interface{}) {
	type OsOnlyInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	i := OsOnlyInfo{
		Name:    si.OperatingSystemName,
		Version: si.OperatingSystemVersion,
	}

	return i
}
