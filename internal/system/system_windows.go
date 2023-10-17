package system

import (
	"errors"
)

// GetSystemProductName returns the product name of the system
func GetSystemProductName() (product string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}

// GetSystemUUID returns the UUID of the system
func GetSystemUUID() (uuid string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}

// GetOSName returns the name of the operating system
func GetOSName() (os string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}

// GetOSVersion returns the version of the operating system
func GetOSVersion() (version string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}
