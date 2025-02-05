package plugin

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os/exec"
	"strings"
)

type NewADOrganizationalUnit struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ADOrganizationalUnitName struct {
	Name []string `json:"name"`
}

type ADOrganizationalUnit struct {
	Name              string
	DistinguishedName string
	Description       string
	ObjectGUID        string
}

// ListADOrganizationalUnits uses PowerShell to fetch all Active Directory organizational units and returns them as a slice of ADOrganizationalUnit structs.
func ListADOrganizationalUnits() ([]ADOrganizationalUnit, error) {
	return listADOrganizationalUnits("'*'")
}

// Get a specific AD organizational unit by Name
func ListADOrganizationalUnit(name string) (*ADOrganizationalUnit, error) {
	filter := fmt.Sprintf("'Name -eq \"%s\"'", name)
	ous, err := listADOrganizationalUnits(filter)
	if err != nil {
		return nil, err
	}
	if len(ous) == 0 {
		return nil, fmt.Errorf("organizational unit not found")
	}
	return &ous[0], nil
}

func listADOrganizationalUnits(filter string) ([]ADOrganizationalUnit, error) {
	// PowerShell script to get specified properties of AD organizational units
	psScript := fmt.Sprintf(`
$filter = %s
Get-ADOrganizationalUnit -Filter $filter -Properties Name, DistinguishedName, Description, ObjectGUID |
Select-Object Name, DistinguishedName, Description, ObjectGUID | 
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

	// Process the output to return a slice of ADOrganizationalUnit
	reader := csv.NewReader(strings.NewReader(out.String()))
	reader.Comma = ',' // Default but re-set for clarity
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV output: %s", err)
	}

	ous := make([]ADOrganizationalUnit, 0, len(records)-1) // Exclude header
	for i, record := range records {
		if i == 0 { // skip the header
			continue
		}
		if len(record) >= 4 {
			ous = append(ous, ADOrganizationalUnit{
				Name:              record[0],
				DistinguishedName: record[1],
				Description:       record[2],
				ObjectGUID:        record[3],
			})
		}
	}

	return ous, nil
}

func CreateADOrganizationalUnit(ou *NewADOrganizationalUnit) (*ADOrganizationalUnit, error) {
	// PowerShell script to create a new AD organizational unit
	psScript := fmt.Sprintf(`
$name = "%s"
$path = "%s"
New-ADOrganizationalUnit -Name $name -Path $path
`, ou.Name, ou.Path)

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

	// Check for errors in the output
	if stderr.Len() > 0 {
		return nil, fmt.Errorf("PowerShell error: %s", stderr.String())
	}

	return ListADOrganizationalUnit(ou.Name)
}

func DeleteADOrganizationalUnit(dn string) error {
	// PowerShell script to delete an AD organizational unit
	psScript := fmt.Sprintf(`
$dn = "%s"
Remove-ADOrganizationalUnit -Identity $dn -Confirm:$false
`, dn)

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
		return fmt.Errorf("failed to execute PowerShell script: %s, with error: %s", stderr.String(), err)
	}

	// Check for errors in the output
	if stderr.Len() > 0 {
		return fmt.Errorf("PowerShell error: %s", stderr.String())
	}

	return nil
}
