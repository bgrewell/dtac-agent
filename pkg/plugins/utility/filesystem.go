package utility

import (
	sharedutil "github.com/bgrewell/dtac-agent/pkg/shared/utility"
)

// FindPlugins returns a list of plugins that match the given pattern
func FindPlugins(root, pattern string) ([]string, error) {
	return sharedutil.Find(root, pattern)
}
