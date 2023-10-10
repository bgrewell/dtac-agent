package system

import (
	"bufio"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"os"
	"strings"
)

func GetSystemProductName() (product string, err error) {
	command := "dmidecode -s system-product-name"
	return helpers.RunAsUser(command, "root")
}

func GetOSName() (os string, err error) {
	return getLinuxInfo("ID")
}

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
			return strings.TrimPrefix(line, fmt.Sprintf("%s=", id)), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("distribution %s not found", id)
}
