package utility

import (
	"errors"
	"os"
	"syscall"
)

func IsOnlyWritableByUserOrRoot(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false, err
	}

	mode := fileInfo.Mode()
	uid := os.Getuid()

	// If running as root, check that file is not writable by group or others
	if uid == 0 {
		return mode&0222 != 0 && mode&0022 == 0, nil
	}

	// For non-root users, check if writable by owner, current user is owner, and not writable by group or others
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return false, errors.New("got unexpected type for file info")
	}
	return mode&0200 != 0 && stat.Uid == uint32(uid) && mode&0022 == 0, nil
}
