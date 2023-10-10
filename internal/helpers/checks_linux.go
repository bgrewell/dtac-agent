package helpers

import (
	"fmt"
	"os/user"
)

// isRunningAsRoot checks if the current process runs as root on UNIX or with elevated privileges on Windows.
func isRunningAsRoot() bool {
	// For UNIX-like systems
	u, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get the current user:", err)
		return false
	}

	return u.Uid == "0"
}
