package diag

import (
	"errors"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
)

// AgentRunningAsUser returns the name and group of the user the agent is running as
func AgentRunningAsUser() (currentUser *types.UserGroup, err error) {
	return nil, errors.New("this function has not been implemented for Windows")
}
