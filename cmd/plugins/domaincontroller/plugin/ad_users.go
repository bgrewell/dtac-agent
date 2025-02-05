package plugin

import (
	"fmt"
	"strconv"
	"strings"
)

type ADUserParam struct {
	Username []string `json:"username"`
}

type ADUserSSHKey struct {
	Username string `json:"username"`
	Key      string `json:"key"`
}

type ADUserSSHKeyImport struct {
	Username       string `json:"username"`
	GitHubUsername string `json:"githubusername"`
}

// ADUser represents an Active Directory user with their relevant properties.
type ADUser struct {
	DisplayName       string
	SamAccountName    string
	UserPrincipalName string
	EmailAddress      string
	AccountDisabled   bool
	SSHPublicKeys     []string
}

type NewADUserObj struct {
	Enabled           bool     `json:"enabled"`
	Username          string   `json:"username"`
	Password          string   `json:"password"`
	UserPrincipalName string   `json:"user_principal_name"`
	DomainName        string   `json:"domain_name"`
	GivenName         string   `json:"given_name"`
	Surname           string   `json:"surname"`
	Company           string   `json:"company"`
	Description       string   `json:"description"`
	DisplayName       string   `json:"display_name"`
	EmailAddress      string   `json:"email_address"`
	Path              string   `json:"path"`
	Shell             string   `json:"shell"`
	SshKeys           []string `json:"ssh_keys"`
}

type UpdateADUserObj struct {
	Enabled           bool     `json:"enabled,omitempty"`
	Username          string   `json:"username"`
	Password          string   `json:"password,omitempty"`
	UserPrincipalName string   `json:"user_principal_name,omitempty"`
	DomainName        string   `json:"domain_name,omitempty"`
	GivenName         string   `json:"given_name,omitempty"`
	Surname           string   `json:"surname,omitempty"`
	Company           string   `json:"company,omitempty"`
	Description       string   `json:"description,omitempty"`
	DisplayName       string   `json:"display_name,omitempty"`
	EmailAddress      string   `json:"email_address,omitempty"`
	Path              string   `json:"path,omitempty"`
	Shell             string   `json:"shell,omitempty"`
	SshKeys           []string `json:"ssh_keys,omitempty"`
}

// ListADUsers returns a list of all AD users.
func ListADUsers() ([]ADUser, error) {
	// Setup script parameters
	params := map[string]string{
		"Filter": "*",
	}
	return listADUsers(params)
}

// ListADUser returns a single AD user given their SamAccountName.
func ListADUser(username string) (*ADUser, error) {
	params := map[string]string{
		"SamAccountName": username,
	}
	users, err := listADUsers(params)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user not found: %s", username)
	}

	return &users[0], nil
}

func CreateADUser(user *NewADUserObj) error {

	params := convertNewADUserToParams(*user)
	params["Name"] = fmt.Sprintf("%s %s", user.GivenName, user.Surname)
	params["Initials"] = user.GivenName[:1] + user.Surname[:1]

	err := executePowerShell(CreateUserScript, params)
	if err != nil {
		return err
	}

	// Add any ssh keys to the users profile
	for _, key := range user.SshKeys {
		keyParams := map[string]string{
			"Username": user.Username,
			"Key":      key,
		}
		err := executePowerShell(CreateSSHKeyScript, keyParams)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteADUser deletes an AD user given their SamAccountName
func DeleteADUser(samAccountName string) error {

	params := map[string]string{
		"SamAccountName": samAccountName,
	}

	err := executePowerShell(DeleteUserScript, params)
	if err != nil {
		return err
	}

	return nil
}

func UpdateADUser(user *NewADUserObj) error {

	params := convertNewADUserToParams(*user)
	err := executePowerShell(ModifyUserScript, params)
	if err != nil {
		return err
	}

	return nil

}

// TODO: ---- EVERYTHING BELOW THIS LINE IS OLD AND NEEDS REWORK

// Helper function to generate the SSH key addition script
func generateSSHKeyScript(username string, sshKeys []string) string {
	if len(sshKeys) == 0 {
		return ""
	}

	var sb strings.Builder
	for _, key := range sshKeys {
		sb.WriteString(fmt.Sprintf(`Set-ADUser -Identity "%s" -Add @{sshPublicKey="%s"}\n`, username, key))
	}
	return sb.String()
}

// TODO -- END OF LIST

// listADUsers returns a list of AD users based on the provided parameters.
func listADUsers(params map[string]string) ([]ADUser, error) {

	records, err := executePowerShellAndParseCsv(ListUserScript, params)
	if err != nil {
		return nil, err
	}

	return parseADUsers(records)
}

func parseADUsers(records [][]string) ([]ADUser, error) {
	if len(records) < 1 {
		return nil, fmt.Errorf("no records found")
	}

	users := make([]ADUser, 0, len(records)-1) // Exclude header
	for i, record := range records {
		if i == 0 { // skip the header
			continue
		}
		if len(record) >= 6 {
			accountEnabled, _ := strconv.ParseBool(record[4])
			sshKeys := []string{}
			if record[5] != "" {
				sshKeys = strings.Split(record[5], ",")
			}
			users = append(users, ADUser{
				DisplayName:       record[0],
				SamAccountName:    record[1],
				UserPrincipalName: record[2],
				EmailAddress:      record[3],
				AccountDisabled:   !accountEnabled,
				SSHPublicKeys:     sshKeys,
			})
		}
	}

	return users, nil
}

func convertNewADUserToParams(user NewADUserObj) map[string]string {
	paramObject := make(map[string]string)

	if user.Username != "" {
		paramObject["SamAccountName"] = user.Username
	}
	if user.Password != "" {
		paramObject["Password"] = user.Password
	}
	if user.UserPrincipalName != "" {
		paramObject["UserPrincipalName"] = user.UserPrincipalName
	}
	if user.GivenName != "" {
		paramObject["GivenName"] = user.GivenName
	}
	if user.Surname != "" {
		paramObject["Surname"] = user.Surname
	}
	if user.Description != "" {
		paramObject["Description"] = user.Description
	}
	if user.EmailAddress != "" {
		paramObject["EmailAddress"] = user.EmailAddress
	}
	if user.Path != "" {
		paramObject["Path"] = user.Path
	}
	// Add any other optional fields here, making sure to check if they are empty

	return paramObject
}
