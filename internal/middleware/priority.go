package middleware

type MiddlewarePriority int

const (
	PriorityAuthentication MiddlewarePriority = 0
	PriorityAuthorization  MiddlewarePriority = 1
	PriorityHigh           MiddlewarePriority = 10
	PriorityMedium         MiddlewarePriority = 100
	PriorityLow            MiddlewarePriority = 200
)
