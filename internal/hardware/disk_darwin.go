package hardware

import "errors"

func NewDiskDetails(name string, size string, model string) *DiskDetails {
	dd := DiskDetails{
		Name:  name,
		Size:  size,
		Model: model,
	}
	return &dd
}

func GetPhysicalDisks() ([]*DiskDetails, error) {
	return nil, errors.New("not implemented")
}
