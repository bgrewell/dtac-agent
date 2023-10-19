package diag

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"os"
	"os/user"
	"syscall"
)

// AgentRunningAsUser returns the name and group of the user the agent is running as
func AgentRunningAsUser() (currentUser *types.UserGroup, err error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	// get the path of the currently running binary
	binaryPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	// get the FileInfo struct describing the binary
	fileInfo, err := os.Stat(binaryPath)
	if err != nil {
		return nil, err
	}

	// get the uid and gid
	gid := fmt.Sprint(fileInfo.Sys().(*syscall.Stat_t).Gid)

	// look up the group name
	g, err := user.LookupGroupId(gid)
	if err != nil {
		return nil, err
	}

	ug := types.UserGroup{
		User:  u.Username,
		Group: g.Name,
	}
	return &ug, nil
}
