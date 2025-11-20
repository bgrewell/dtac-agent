package utility

import (
	sharedutil "github.com/bgrewell/dtac-agent/pkg/shared/utility"
)

// GetUnusedTCPPort returns an unused TCP port
func GetUnusedTCPPort() (int, error) {
	return sharedutil.GetUnusedTCPPort()
}
