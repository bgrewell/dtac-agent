<#
    .SYNOPSIS
        Set-PlatformUserSSHKey - Add an SSH public key to a specified user in Active Directory.

    .DESCRIPTION
        This script adds an SSH public key to the 'sshPublicKey' attribute for a specified user in Active Directory.

    .PARAMETER Username
        The username of the user to whom the SSH public key needs to be added. This parameter is required.

    .PARAMETER Key
        The SSH public key that needs to be added to the specified user. This parameter is required.

    .EXAMPLE
        .\Set-PlatformUserSSHKey.ps1 -Username jdoe -Key "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA..."

    .NOTES
        FileName:        Set-PlatformUserSSHKey.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>


param(
    [Parameter(Mandatory=$true)]
    [string]$Username,
    [Parameter(Mandatory=$true)]
    [string]$Key
)

function Set-PlatformUserSSHKey {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Username,
        [Parameter(Mandatory=$true)]
        [string]$Key
    )

    # Check if the user exists in Active Directory
    $ADUser = Get-ADUser -Filter { SamAccountName -eq $Username } -Properties sshPublicKey

    if ($ADUser) {
        # Add the public key to the 'sshPublicKey' attribute
        try {
            If ($ADUser.sshPublicKey -contains $Key) {
                Write-Error "The key is already present. No changes made."
                exit 1
            } else {
                Set-ADUser -Identity $ADUser -Add @{sshPublicKey=@($Key)}
                Write-Host "The key has been added to the user's sshPublicKey attribute." -ForegroundColor Green
            }
        } catch {
            Write-Error "Failed to add the key to the user's sshPublicKey attribute: $_"
            exit 1
        }

    } else {
        Write-Error "User ($Username) not found in Active Directory."
        exit 1
    }
}

# To call the function from the script
Set-PlatformUserSSHKey -Username $Username -Key "$Key"