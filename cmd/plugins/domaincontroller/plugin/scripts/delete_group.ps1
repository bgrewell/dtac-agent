<#
    .SYNOPSIS
        Remove-PlatformGroup - Remove a group from Active Directory.

    .DESCRIPTION
        This script removes a group from Active Directory based on the specified SamAccountName.

    .PARAMETER SamAccountName
        The SamAccountName of the group to be removed. This parameter is required.

    .EXAMPLE
        .\Remove-PlatformGroup.ps1 -SamAccountName "GroupName"

    .NOTES
        FileName:        Remove-PlatformGroup.ps1
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

function Remove-PlatformGroup {
    param (
        [Parameter(Mandatory=$true)][string]$SamAccountName
    )

    # Remove the group
    Remove-ADGroup -Identity $SamAccountName -Confirm:$false
}

# To call the function from the script
Remove-PlatformGroup -SamAccountName $SamAccountName
