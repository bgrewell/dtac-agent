<#
    .SYNOPSIS
        Set-PlatformGroup - Modify a group in Active Directory.

    .DESCRIPTION
        This script modifies a group in Active Directory based on the specified parameters.

    .PARAMETER SamAccountName
        The SamAccountName of the group to be modified. This parameter is required.

    .PARAMETER Description
        The new description of the group.

    .PARAMETER GroupScope
        The new scope of the group (e.g., Global, Universal, DomainLocal).

    .PARAMETER ManagedBy
        The distinguished name (DN) of the user or group that manages this group.

    .EXAMPLE
        .\Set-PlatformGroup.ps1 -SamAccountName "GroupName" -Description "New Description" -GroupScope "Universal" -ManagedBy "CN=Manager,OU=Users,DC=example,DC=com"

    .NOTES
        FileName:        Set-PlatformGroup.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$SamAccountName,
    [string]$Description,
    [string]$GroupScope,
    [string]$ManagedBy
)

function Set-PlatformGroup {
    param (
        [Parameter(Mandatory=$true)][string]$SamAccountName,
        [string]$Description,
        [string]$GroupScope,
        [string]$ManagedBy
    )

    # Modify the group
    Set-ADGroup -Identity $SamAccountName `
        -Description $Description `
        -GroupScope $GroupScope `
        -ManagedBy $ManagedBy
}

# To call the function from the script
Set-PlatformGroup -SamAccountName $SamAccountName -Description $Description -GroupScope $GroupScope -ManagedBy $ManagedBy
