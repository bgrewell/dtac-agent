<#
    .SYNOPSIS
        Remove-PlatformUserSSHKey - Remove an SSH public key from a specified user in Active Directory.

    .DESCRIPTION
        This script removes an SSH public key from the 'sshPublicKey' attribute for a specified user in Active Directory.

    .PARAMETER Username
        The username of the user from whom the SSH public key needs to be removed. This parameter is required.

    .PARAMETER Key
        The SSH public key that needs to be removed from the specified user. This parameter is required.

    .EXAMPLE
        .\Remove-PlatformUserSSHKey.ps1 -Username jdoe -Key "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAA..."

    .NOTES
        FileName:        Remove-PlatformUserSSHKey.ps1
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

function Remove-PlatformUserSSHKey {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Username,
        [Parameter(Mandatory=$true)]
        [string]$Key
    )

    # Import active directory module for running AD cmdlets
    Import-Module ActiveDirectory

    # Check if the user exists in Active Directory
    $ADUser = Get-ADUser -Filter { SamAccountName -eq $Username } -Properties sshPublicKey

    if ($ADUser) {
        # Remove the public key from the 'sshPublicKey' attribute
        try {
            if ($ADUser.sshPublicKey -contains $Key) {
                Set-ADUser -Identity $ADUser -Remove @{sshPublicKey=@($Key)}
                Write-Host "The key has been removed from the user's sshPublicKey attribute." -ForegroundColor Green
            } else {
                Write-Error "The key is not present. No changes made."
                exit 1
            }
        } catch {
            Write-Error "Failed to remove the key from the user's sshPublicKey attribute: $_"
            exit 1
        }

        # Display user sshPublicKey values
        $ADUser = Get-ADUser -Identity $Username -Properties sshPublicKey
        Write-Host "All User SSH Public Keys  :  "
        foreach ($SshKey in $ADUser.sshPublicKey) {
            Write-Host
            Write-Host $SshKey -ForegroundColor Blue
        }
        Write-Host ""
    } else {
        Write-Error "User ($Username) not found in Active Directory."
        exit 1
    }
}

# To call the function from the script
Remove-PlatformUserSSHKey -Username $Username -Key $Key
