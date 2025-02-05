<#
    .SYNOPSIS
        New-PlatformOrganizationalUnit - Create a new organizational unit (OU) in Active Directory.

    .DESCRIPTION
        This script creates a new organizational unit (OU) in Active Directory with the specified parameters.

    .PARAMETER Name
        The name of the new organizational unit. This parameter is required.

    .PARAMETER Path
        The path in the directory where the new organizational unit will be created. This parameter is required.

    .EXAMPLE
        .\New-PlatformOrganizationalUnit.ps1 -Name "NewOU" -Path "OU=Departments,DC=example,DC=com"

    .NOTES
        FileName:        New-PlatformOrganizationalUnit.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$Name,
    [Parameter(Mandatory=$true)][string]$Path
)

function New-PlatformOrganizationalUnit {
    param (
        [Parameter(Mandatory=$true)][string]$Name,
        [Parameter(Mandatory=$true)][string]$Path
    )

    # Create the new organizational unit
    New-ADOrganizationalUnit -Name $Name -Path $Path
}

# To call the function from the script
New-PlatformOrganizationalUnit -Name $Name -Path $Path
