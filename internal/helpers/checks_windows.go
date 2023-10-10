package helpers

import (
	"os"
)

// isRunningAsRoot checks if the current process runs as root on UNIX or with elevated privileges on Windows.
func isRunningAsRoot() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}
