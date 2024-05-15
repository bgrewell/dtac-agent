package system

import (
	"errors"
	"fmt"
	"github.com/StackExchange/wmi"
)

type Win32_ComputerSystemProduct struct {
	UUID string
}

// GetSystemProductName returns the product name of the system
func GetSystemProductName() (product string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}

// GetSystemUUID returns the UUID of the system
func GetSystemUUID() (uuid string, err error) {
	var dst []Win32_ComputerSystemProduct
	query := "SELECT UUID FROM Win32_ComputerSystemProduct"

	// Perform the WMI query
	if err := wmi.Query(query, &dst); err != nil {
		return "", fmt.Errorf("WMI Query failed: %s", err)
	}

	// Output the UUID
	if len(dst) > 0 {
		return dst[0].UUID, nil
	} else {
		return "", errors.New("No UUID found")
	}
}

// GetOSName returns the name of the operating system
func GetOSName() (os string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}

// GetOSVersion returns the version of the operating system
func GetOSVersion() (version string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}
