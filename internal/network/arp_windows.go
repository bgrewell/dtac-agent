package network

import (
	"os/exec"
	"strings"
)

// GetArpTable returns the ARP table as a slice of ArpEntry structs.
func GetArpTable() ([]ArpEntry, error) {
	cmd := exec.Command("arp", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var arpEntries []ArpEntry

	// Start from 1 to skip the header
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			entry := ArpEntry{
				IPAddress: fields[0],
				HWType:    fields[1],
				Flags:     fields[2],
				HWAddress: fields[3],
				Mask:      fields[4],
				Iface:     fields[5],
			}
			arpEntries = append(arpEntries, entry)
		}
	}

	return arpEntries, nil
}
