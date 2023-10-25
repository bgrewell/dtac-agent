package endpoint

// Action is an enumeration for the actions that an endpoint can execute
type Action string

const (
	// ActionCreate is the action for creating a resource, e.g. POST in REST
	ActionCreate Action = "create"
	// ActionRead is the action for reading a resource, e.g. GET in REST
	ActionRead = "read"
	// ActionWrite is the action for writing a resource, e.g. PUT in REST
	ActionWrite = "write"
	// ActionDelete is the action for deleting a resource, e.g. DELETE in REST
	ActionDelete = "delete"
)
