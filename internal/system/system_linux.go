package system

import (
	"bufio"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"os"
	"strings"
)

// GetSystemProductName returns the product name of the system
func GetSystemProductName() (product string, err error) {
	command := "dmidecode -s system-product-name"
	return helpers.RunAsUser(command, "root")
}

// GetSystemUUID returns the UUID of the system
func GetSystemUUID() (uuid string, err error) {
	command := "dmidecode -s system-uuid"
	return helpers.RunAsUser(command, "root")
}

// GetOSName returns the name of the operating system
func GetOSName() (os string, err error) {
	return getLinuxInfo("ID")
}

// GetOSVersion returns the version of the operating system
func GetOSVersion() (version string, err error) {
	return getLinuxInfo("VERSION")
}

func getLinuxInfo(id string) (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, fmt.Sprintf("%s=", id)) {
			val := strings.TrimPrefix(line, fmt.Sprintf("%s=", id))
			val = strings.Trim(val, "\"")
			return val, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("distribution %s not found", id)
}
