package helpers

// IsRunningAsRoot returns true if the process is running as root
func IsRunningAsRoot() bool {
	return isRunningAsRoot()
}
