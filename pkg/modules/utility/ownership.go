package utility

import (
	sharedutil "github.com/bgrewell/dtac-agent/pkg/shared/utility"
)

// IsOnlyWritableByUserOrRoot checks if the file is only writable by the user or root.
func IsOnlyWritableByUserOrRoot(filename string) (onlyWritable bool, err error) {
	return sharedutil.IsOnlyWritableByUserOrRoot(filename)
}
