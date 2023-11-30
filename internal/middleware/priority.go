package middleware

// Priority is the priority for middleware and is used to control call chaining
type Priority int

const (
	// PriorityAuthentication is the priority for authentication middleware
	PriorityAuthentication Priority = 0
	// PriorityAuthorization is the priority for authorization middleware
	PriorityAuthorization Priority = 1
	// PriorityHigh is for high priority non-authn/non-authz middleware
	PriorityHigh Priority = 10
	// PriorityMedium is for medium priority non-authn/non-authz middleware
	PriorityMedium Priority = 100
	// PriorityLow is for low priority non-authn/non-authz middleware
	PriorityLow Priority = 200
	// PriorityValidation is for validation middleware which is one of the last to run
	PriorityValidation = 300
)
