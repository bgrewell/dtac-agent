package utility

import (
	pluginutil "github.com/bgrewell/dtac-agent/pkg/plugins/utility"
)

// IsOnlyWritableByUserOrRoot checks if the file is only writable by the user or root.
func IsOnlyWritableByUserOrRoot(filename string) (onlyWritable bool, err error) {
	return pluginutil.IsOnlyWritableByUserOrRoot(filename)
}
