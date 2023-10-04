package diag

import "runtime"

// CurrentMemoryStats returns the current runtime memory statistics
func CurrentMemoryStats() *runtime.MemStats {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return &m
}
