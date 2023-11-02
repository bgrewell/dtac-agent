package endpoint

import (
	"fmt"
	"strings"
)

// Action is an enumeration for the actions that an endpoint can execute
type Action string

// String converts the Action type to a string.
func (a Action) String() string {
	return string(a)
}

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

// ParseAction converts a string to an Action type.
func ParseAction(s string) (Action, error) {
	switch strings.ToLower(s) {
	case "create":
		return ActionCreate, nil
	case "read":
		return ActionRead, nil
	case "write":
		return ActionWrite, nil
	case "delete":
		return ActionDelete, nil
	default:
		return "", fmt.Errorf("invalid action: %s", s)
	}
}
