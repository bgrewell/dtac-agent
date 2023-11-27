package utility

import (
	"golang.org/x/sys/windows"
)

// IsOnlyWritableByUserOrRoot checks if the file is only writable by the user or root.
func IsOnlyWritableByUserOrRoot(filename string) (onlyWritable bool, err error) {
	// Get the Windows SID of the current user
	user, err := windows.CurrentUser()
	if err != nil {
		return false, err
	}
	userSID, err := windows.StringToSid(user.Sid)
	if err != nil {
		return false, err
	}

	// Get the security descriptor of the file
	securityDescriptor, err := windows.GetNamedSecurityInfo(filename, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return false, err
	}
	dacl, _, err := securityDescriptor.DACL()
	if err != nil {
		return false, err
	}

	// Check the DACL for user-specific write permissions
	for _, ace := range dacl {
		if ace.Type == windows.ACCESS_ALLOWED_ACE_TYPE {
			if ace.Mask&windows.GENERIC_WRITE == windows.GENERIC_WRITE {
				sid, err := ace.Sid()
				if err != nil {
					return false, err
				}
				if !windows.EqualSid(sid, userSID) {
					return false, nil
				}
			}
		}
	}
	return true, nil
}
