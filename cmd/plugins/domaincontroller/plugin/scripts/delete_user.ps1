<#
    .SYNOPSIS
        Remove-PlatformUser - Remove a user from Active Directory.

    .DESCRIPTION
        This script removes a user from Active Directory based on the specified SamAccountName.

    .PARAMETER SamAccountName
        The SamAccountName of the user to be removed. This parameter is required.

    .EXAMPLE
        .\Remove-PlatformUser.ps1 -SamAccountName "jdoe"

    .NOTES
        FileName:        Remove-PlatformUser.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$SamAccountName
)

function Remove-PlatformUser {
    param (
        [Parameter(Mandatory=$true)][string]$SamAccountName
    )

    # Remove the user
    Remove-ADUser -Identity $SamAccountName -Confirm:$false
}

# To call the function from the script
Remove-PlatformUser -SamAccountName $SamAccountName
