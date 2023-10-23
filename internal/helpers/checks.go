package helpers

import (
	"encoding/json"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"io"
	"net/http"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// IsRunningAsRoot returns true if the process is running as root
func IsRunningAsRoot() bool {
	return isRunningAsRoot()
}

// CheckUserGroup checks whether the current user is in the specified group.
func CheckUserGroup(currentUser *user.User, group string) bool {
	userGID, err := strconv.Atoi(currentUser.Gid)
	if err != nil {
		return false
	}

	groupUser, err := user.LookupGroup(group)
	if err != nil {
		return false
	}

	groupGID, err := strconv.Atoi(groupUser.Gid)
	if err != nil {
		return false
	}

	return isInGroup(userGID, groupGID)
}

// CheckUser checks if the current user matches the criteria returned by the API.
func CheckUser(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var userGroup types.UserGroup
	err = json.Unmarshal(body, &userGroup)
	if err != nil {
		return false, err
	}

	currentUser, err := user.Current()
	if err != nil {
		return false, err
	}

	isUser := currentUser.Username == userGroup.User
	isGroup := CheckUserGroup(currentUser, userGroup.Group) // assuming you have this function implemented
	isRoot := currentUser.Username == "root"

	return isUser || isGroup || isRoot, nil
}

// CanRead checks whether a file is readable.
func CanRead(path string) bool {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		// unable to open the file for reading -> return false
		return false
	}
	defer file.Close()

	return true
}

func isInGroup(userGID, groupGID int) bool {
	groups, err := syscall.Getgroups()
	if err != nil {
		return false
	}

	for _, g := range groups {
		if g == userGID || g == groupGID {
			return true
		}
	}
	return false
}
