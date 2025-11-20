package utility

import (
	pluginutil "github.com/bgrewell/dtac-agent/pkg/plugins/utility"
)

// GetUnusedTCPPort returns an unused TCP port
func GetUnusedTCPPort() (int, error) {
	return pluginutil.GetUnusedTCPPort()
}
