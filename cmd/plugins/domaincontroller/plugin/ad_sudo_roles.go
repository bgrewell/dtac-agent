package plugin

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"
)

type ADSudoRoleName struct {
	Name []string `json:"name"`
}

type NewSudoRole struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SudoHost    string `json:"sudo_host"`
	SudoUser    string `json:"sudo_user"`
	SudoCommand string `json:"sudo_command"`
}

type DeleteSudoRole struct {
	Name string `json:"name"`
}

// ADSudoRole represents an Active Directory sudoRole object with relevant properties.
type ADSudoRole struct {
	Name              string
	DistinguishedName string
	CN                string
	SudoCommand       []string
	SudoHost          []string
	SudoUser          []string
}

// ListADSudoRoles uses PowerShell to fetch all Active Directory sudoRole objects and returns them as a slice of ADSudoRole structs.
func ListADSudoRoles() ([]ADSudoRole, error) {
	return listADSudoRoles("'ObjectClass -eq \"sudoRole\"'")
}

func ListADSudoRole(name string) (*ADSudoRole, error) {
	filter := fmt.Sprintf("'ObjectClass -eq \"sudoRole\" -and Name -eq \"%s\"'", name)
	sudoRoles, err := listADSudoRoles(filter)
	if err != nil {
		return nil, err
	}
	if len(sudoRoles) == 0 {
		return nil, fmt.Errorf("sudo role not found")
	}
	return &sudoRoles[0], nil
}

func listADSudoRoles(filter string) ([]ADSudoRole, error) {
	// PowerShell script to get specified properties of AD sudoRole objects
	psScript := fmt.Sprintf(`
$filter = %s
Get-ADObject -Filter $filter -Properties Name, DistinguishedName, CN, sudoCommand, sudoHost, sudoUser |
Select-Object Name, DistinguishedName, CN, 
@{Name='SudoCommand';Expression={($_.sudoCommand -join ';')}},
@{Name='SudoHost';Expression={($_.sudoHost -join ';')}},
@{Name='SudoUser';Expression={($_.sudoUser -join ';')}} |
ConvertTo-Csv -NoTypeInformation
`, filter)

	// Set up the command to run the PowerShell script
	cmd := exec.Command("powershell", "-Command", psScript)

	// Capture the output and errors from the command
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell script: %s, with error: %s", stderr.String(), err)
	}

	// Process the output to return a slice of ADSudoRole
	reader := csv.NewReader(strings.NewReader(out.String()))
	reader.Comma = ',' // Default but re-set for clarity
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV output: %s", err)
	}

	sudoRoles := make([]ADSudoRole, 0, len(records)-1) // Exclude header
	for i, record := range records {
		if i == 0 { // skip the header
			continue
		}
		if len(record) >= 6 {
			sudoRoles = append(sudoRoles, ADSudoRole{
				Name:              record[0],
				DistinguishedName: record[1],
				CN:                record[2],
				SudoCommand:       strings.Split(record[3], ";"),
				SudoHost:          strings.Split(record[4], ";"),
				SudoUser:          strings.Split(record[5], ";"),
			})
		}
	}

	return sudoRoles, nil
}
