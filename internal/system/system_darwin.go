package system

import (
	"errors"
)

func GetSystemProductName() (product string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}

func GetOSName() (os string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}

func GetOSVersion() (version string, err error) {
	return "", errors.New("this function has not been implemented for this OS")
}
