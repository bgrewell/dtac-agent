package utility

import (
	sharedutil "github.com/bgrewell/dtac-agent/pkg/shared/utility"
)

// FindModules returns a list of modules that match the given pattern
func FindModules(root, pattern string) ([]string, error) {
	return sharedutil.Find(root, pattern)
}
