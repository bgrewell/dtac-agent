package plugin

import (
	"fmt"
)

type NewADGroup struct {
	Name string `json:"name"`
}

type ADGroupName struct {
	Name []string `json:"name"`
}

// ADGroup represents an Active Directory group with their relevant properties.
type ADGroup struct {
	Name              string
	SamAccountName    string
	DistinguishedName string
	GroupCategory     string
	GroupScope        string
	SID               string
}

type DeleteADGroup struct {
	Name string `json:"name"`
}

// ListADGroups uses PowerShell to fetch all Active Directory groups and returns them as a slice of ADGroup structs.
func ListADGroups() ([]ADGroup, error) {
	params := map[string]string{
		"Filter": "*",
	}
	return listADGroups(params)
}

func ListADGroup(name string) (*ADGroup, error) {
	filter := fmt.Sprintf("'Name -eq \"%s\"'", name)
	params := map[string]string{
		"Filter": filter,
	}
	groups, err := listADGroups(params)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return nil, fmt.Errorf("group not found")
	}
	return &groups[0], nil
}

func UpdateADGroup()

func listADGroups(params map[string]string) ([]ADGroup, error) {
	records, err := executePowerShellAndParseCsv(ListUserScript, params)
	if err != nil {
		return nil, err
	}

	return parseADGroups(records)
}
func parseADGroups(records [][]string) ([]ADGroup, error) {

	if len(records) < 1 {
		return nil, fmt.Errorf("no records found")
	}

	groups := make([]ADGroup, 0, len(records)-1) // Exclude header
	for i, record := range records {
		if i == 0 { // skip the header
			continue
		}
		if len(record) >= 6 {
			groups = append(groups, ADGroup{
				Name:              record[0],
				SamAccountName:    record[1],
				DistinguishedName: record[2],
				GroupCategory:     record[3],
				GroupScope:        record[4],
				SID:               record[5],
			})
		}
	}

	return groups, nil
}
