<#
    .SYNOPSIS
        Remove-PlatformOrganizationalUnit - Remove an organizational unit (OU) from Active Directory.

    .DESCRIPTION
        This script removes an organizational unit (OU) from Active Directory based on the specified distinguished name (DN).

    .PARAMETER dn
        The distinguished name (DN) of the organizational unit to be removed. This parameter is required.

    .EXAMPLE
        .\Remove-PlatformOrganizationalUnit.ps1 -dn "OU=Departments,DC=example,DC=com"

    .NOTES
        FileName:        Remove-PlatformOrganizationalUnit.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$dn
)

function Remove-PlatformOrganizationalUnit {
    param (
        [Parameter(Mandatory=$true)][string]$dn
    )

    # Remove the organizational unit
    Remove-ADOrganizationalUnit -Identity $dn -Confirm:$false
}

# To call the function from the script
Remove-PlatformOrganizationalUnit -dn $dn
