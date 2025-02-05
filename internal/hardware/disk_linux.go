package hardware

import (
	execute "github.com/bgrewell/go-execute/v2"
	"github.com/shirou/gopsutil/disk"
	"strings"
)

// NewDiskDetails creates a new DiskDetails struct
func NewDiskDetails(name string, size string, model string) *DiskDetails {
	dd := DiskDetails{
		Name:   name,
		Size:   size,
		Model:  model,
		Serial: disk.GetDiskSerialNumber(name),
		Label:  disk.GetLabel(name),
	}
	return &dd
}

// GetPhysicalDisks returns the physical disks
func GetPhysicalDisks() ([]*DiskDetails, error) {
	ex := execute.NewExecutor()
	stdout, _, err := ex.ExecuteSeparate("lsblk -o NAME,TYPE,SIZE,MODEL")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(stdout, "\n")
	var disks []*DiskDetails
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 4 && fields[1] == "disk" {
			model := strings.Join(fields[3:], " ")
			disk := NewDiskDetails(fields[0], fields[2], model)
			disks = append(disks, disk)
		}
	}

	return disks, nil
}
