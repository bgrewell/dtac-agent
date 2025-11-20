package utility

import (
	sharedutil "github.com/bgrewell/dtac-agent/pkg/shared/utility"
)

// IsOnlyWritableByUserOrRoot checks if the file is only writable by the user or root.
func IsOnlyWritableByUserOrRoot(filename string) (bool, error) {
	return sharedutil.IsOnlyWritableByUserOrRoot(filename)
}
