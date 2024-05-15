package diag

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"os/user"
)

// AgentRunningAsUser returns the name and group of the user the agent is running as
func AgentRunningAsUser() (currentUser *types.UserGroup, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	ug := &types.UserGroup{
		User:  u.Username,
		Group: u.Gid,
	}
	return ug, nil
}
