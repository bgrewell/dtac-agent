package hardware

import "errors"

// NewDiskDetails creates a new DiskDetails struct
func NewDiskDetails(name string, size string, model string) *DiskDetails {
	dd := DiskDetails{
		Name:  name,
		Size:  size,
		Model: model,
	}
	return &dd
}

// GetPhysicalDisks returns the physical disks
func GetPhysicalDisks() ([]*DiskDetails, error) {
	return nil, errors.New("not implemented")
}
