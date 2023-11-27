package utility

import "errors"

// IsOnlyWritableByUserOrRoot checks if the file is only writable by the user or root.
func IsOnlyWritableByUserOrRoot(filename string) (bool, error) {
	return false, errors.New("this method has not been implemented for Windows yet")
	//// Get the current user
	//user, err := user.Current()
	//if err != nil {
	//	return false, err
	//}
	//userSID, err := windows.StringToSid(user.Uid)
	//if err != nil {
	//	return false, err
	//}
	//
	//// Get the security descriptor of the file
	//securityDescriptor, err := windows.GetNamedSecurityInfo(filename, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	//if err != nil {
	//	return false, err
	//}
	//dacl, _, err := securityDescriptor.DACL()
	//if err != nil {
	//	return false, err
	//}
	//
	//// Convert the DACL to explicit entries
	//var aceCount uint32
	//var aclSizeInfo windows.AclSizeInformation
	//err = windows.GetAclInformation(dacl, &aclSizeInfo, windows.AclSizeInformationEnum)
	//if err != nil {
	//	return false, err
	//}
	//aceCount = aclSizeInfo.AceCount
	//
	//// Iterate through the ACEs in the DACL
	//for i := uint32(0); i < aceCount; i++ {
	//	var ace *windows.AceHeader
	//	err = windows.GetAce(dacl, i, &ace)
	//	if err != nil {
	//		return false, err
	//	}
	//
	//	if ace.AceType == windows.ACCESS_ALLOWED_ACE_TYPE {
	//		// Additional logic to check for write permissions...
	//	}
	//}
	//
	//// Additional implementation...
	//return true, nil
}
