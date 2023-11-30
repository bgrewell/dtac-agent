package types

// ContextKey is a enum to help manage keys used in passing values in context
type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

const (
	// ContextExecDuration is the key used to store the value of the execution duration
	ContextExecDuration ContextKey = "exec_duration"
	// ContextAuthHeader is the key used to store the value of the auth header
	ContextAuthHeader ContextKey = "auth_header"
	// ContextAuthUser is the key used to store the value of the auth user
	ContextAuthUser ContextKey = "auth_user"
	// ContextAuthRoles is the key used to store the value of the auth roles
	ContextAuthRoles ContextKey = "auth_roles"
	// ContextResourceAction is the key used to store the value of the resource action
	ContextResourceAction ContextKey = "resource_action"
	// ContextResourcePath is the key used to store the value of the resource path
	ContextResourcePath ContextKey = "resource_path"
)
