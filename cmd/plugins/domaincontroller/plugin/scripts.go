package plugin

import (
	_ "embed"
	"encoding/csv"
	"fmt"
	"github.com/bgrewell/go-execute/v2"
	"strings"
)

// === User ===
//
//go:embed scripts/list_user.ps1
var ListUserScript string

//go:embed scripts/create_user.ps1
var CreateUserScript string

//go:embed scripts/delete_user.ps1
var DeleteUserScript string

//go:embed scripts/modify_user.ps1
var ModifyUserScript string

// === Group ===
//
//go:embed scripts/list_group.ps1
var ListGroupScript string

//go:embed scripts/create_group.ps1
var CreateGroupScript string

//go:embed scripts/delete_group.ps1
var DeleteGroupScript string

//go:embed scripts/modify_group.ps1
var ModifyGroupScript string

// === OU ===
//
//go:embed scripts/list_ou.ps1
var ListOUScript string

//go:embed scripts/create_ou.ps1
var CreateOUScript string

//go:embed scripts/delete_ou.ps1
var DeleteOUScript string

//go:embed scripts/modify_ou.ps1
var ModifyOUScript string

// === Sudo Role ===
//
//go:embed scripts/list_sudo_role.ps1
var ListSudoRoleScript string

//go:embed scripts/create_sudo_role.ps1
var CreateSudoRoleScript string

//go:embed scripts/delete_sudo_role.ps1
var DeleteSudoRoleScript string

//go:embed scripts/modify_sudo_role.ps1
var ModifySudoRoleScript string

// === SSH Keys ===
//
//go:embed scripts/create_ssh_key.ps1
var CreateSSHKeyScript string

//go:embed scripts/delete_ssh_key.ps1
var DeleteSSHKeyScript string

//go:embed scripts/import_ssh_keys.ps1
var ImportSSHKeyScript string

func executePowerShell(script string, params map[string]string) (err error) {
	_, stderr, err := executor.ExecuteScriptFromString(execute.ScriptTypePowerShell, script, nil, params)
	if err != nil {
		return fmt.Errorf("failed to execute PowerShell script: %s\n%s", err, stderr)
	} else if stderr != "" {
		return fmt.Errorf("PowerShell script error: %s", stderr)
	}
	return nil
}

func executePowerShellAndParseCsv(script string, params map[string]string) (records [][]string, err error) {
	stdout, stderr, err := executor.ExecuteScriptFromString(execute.ScriptTypePowerShell, script, nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to execute PowerShell script: %s\n%s", err, stderr)
	} else if stderr != "" {
		return nil, fmt.Errorf("PowerShell script error: %s", stderr)
	}

	reader := csv.NewReader(strings.NewReader(stdout))
	reader.Comma = ','
	reader.TrimLeadingSpace = true

	records, err = reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV output: %s", err)
	}

	return records, nil
}
