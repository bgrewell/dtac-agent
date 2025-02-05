<#
    .SYNOPSIS
        Import-PlatformUserSSHKeys - Retrieve SSH public keys from a GitHub user and add them to the specified Active Directory user.

    .DESCRIPTION
        This script retrieves SSH public keys from a GitHub user and adds them to the 'sshPublicKey' attribute of the specified Active Directory user.

    .PARAMETER Username
        The Active Directory username to whom the SSH public keys will be added.

    .PARAMETER GitHubUsername
        The GitHub username from which the SSH public keys will be fetched.

    .EXAMPLE
        .\Import-PlatformUserSSHKeys.ps1 -Username jdoe -GitHubUsername johndoe

    .NOTES
        FileName:        Import-PlatformUserSSHKeys.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)
    History:
        27-Mar-2023: Initial Version
        19-Apr-2023: Made idempotent so it can be ran multiple times on a user to pickup recent changes
    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory = $true)]
    [string]$Username,

    [Parameter(Mandatory = $true)]
    [string]$GitHubUsername
)

function Import-PlatformUserSSHKeys {
    param (
        [Parameter(Mandatory = $true)]
        [string]$Username,

        [Parameter(Mandatory = $true)]
        [string]$GitHubUsername
    )

    # Function to get the public keys from a GitHub user
    function Get-GitHubPublicKeys {
        param(
            [string]$GitHubUsername
        )
        $GitHubApiUrl = "https://api.github.com/users/$GitHubUsername/keys"
        $response = Invoke-WebRequest -Uri $GitHubApiUrl -UseBasicParsing

        if ($response.StatusCode -eq 200) {
            $keys = ($response.Content | ConvertFrom-Json) | Select-Object -ExpandProperty key
            return $keys
        } else {
            Write-Error "Failed to fetch public keys for GitHub user $GitHubUsername"
            exit 1
        }
    }

    # Retrieve the public keys from the user's GitHub account
    $publicKeys = Get-GitHubPublicKeys -GitHubUsername $GitHubUsername

    # Check and add public keys to the 'sshPublicKey' attribute if not already present
    if ($publicKeys.Count -gt 0) {
        $currentUser = Get-ADUser -Identity $Username -Properties sshPublicKey
        foreach ($key in $publicKeys) {
            if (-not ($currentUser.sshPublicKey -contains $key)) {
                try {
                    Set-ADUser -Identity $Username -Add @{sshPublicKey=@($key.ToString())}
                    Write-Host "Added SSH key for user $Username : $key"
                } catch {
                    Write-Error "Failed to add key for user $Username : $_"
                    exit 1
                }
            } else {
                Write-Host "SSH key already exists for user $Username : $key"
            }
        }
    } else {
        Write-Error "No SSH keys found for GitHub user $GitHubUsername."
        exit 1
    }
}

# To call the function from the script
Import-PlatformUserSSHKeys -Username $Username -GitHubUsername $GitHubUsername
