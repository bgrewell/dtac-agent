package network

import (
	"bufio"
	"errors"
	"os/exec"
	"strings"
)

func GetArpTable() ([]ArpEntry, error) {
	cmd := exec.Command("arp", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var entries []ArpEntry

	// Read the header line
	if !scanner.Scan() {
		return nil, errors.New("unexpected arp output format")
	}

	headers := scanner.Text()
	fieldPositions := getFieldPositions(headers)

	for scanner.Scan() {
		line := scanner.Text()
		entry := ArpEntry{
			IPAddress: extractField(line, fieldPositions[0]),
			HWType:    extractField(line, fieldPositions[1]),
			HWAddress: extractField(line, fieldPositions[2]),
			Flags:     extractField(line, fieldPositions[3]),
			Mask:      extractField(line, fieldPositions[4]),
			Iface:     extractField(line, fieldPositions[5]),
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func getFieldPositions(header string) []int {
	var positions []int
	lastChar := ' '
	for i, char := range header {
		if char != ' ' && lastChar == ' ' {
			positions = append(positions, i)
		}
		lastChar = char
	}
	return positions
}

func extractField(line string, start int) string {
	end := start
	for ; end < len(line) && line[end] != ' '; end++ {
	}
	return strings.TrimSpace(line[start:end])
}
