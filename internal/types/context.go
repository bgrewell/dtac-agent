package types

// ContextKey is a enum to help manage keys used in passing values in context
type ContextKey string

const (
	// ContextExecDuration is the key used to store the value of the execution duration
	ContextExecDuration ContextKey = "exec_duration"
)
