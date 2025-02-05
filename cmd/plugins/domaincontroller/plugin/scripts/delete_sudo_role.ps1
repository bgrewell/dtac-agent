<#
    .SYNOPSIS
        Remove-PlatformSudoRole - Remove a sudo role object from Active Directory.

    .DESCRIPTION
        This script removes a sudo role object from Active Directory based on the specified RoleName and Path.

    .PARAMETER RoleName
        The name of the sudo role to be removed. This parameter is required.

    .PARAMETER Path
        The path in the directory where the sudo role is located. This parameter is required.

    .EXAMPLE
        .\Remove-PlatformSudoRole.ps1 -RoleName "AdminRole" -Path "OU=Sudoers,DC=example,DC=com"

    .NOTES
        FileName:        Remove-PlatformSudoRole.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$RoleName,
    [Parameter(Mandatory=$true)][string]$Path
)

function Remove-PlatformSudoRole {
    param (
        [Parameter(Mandatory=$true)][string]$RoleName,
        [Parameter(Mandatory=$true)][string]$Path
    )

    # Get the distinguished name of the sudo role object
    $sudoRoleDN = (Get-ADObject -Filter "Name -eq '$RoleName'" -SearchBase $Path).DistinguishedName

    # Remove the sudo role object
    Remove-ADObject -Identity $sudoRoleDN -Confirm:$false
}

# To call the function from the script
Remove-PlatformSudoRole -RoleName $RoleName -Path $Path
