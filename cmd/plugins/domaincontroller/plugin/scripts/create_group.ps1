<#
    .SYNOPSIS
        New-PlatformGroup - Create a new group in Active Directory.

    .DESCRIPTION
        This script creates a new group in Active Directory with the specified parameters.

    .PARAMETER Name
        The name of the new group. This parameter is required.

    .PARAMETER SamAccountName
        The SamAccountName of the new group. This parameter is required.

    .PARAMETER GroupScope
        The scope of the new group (e.g., Global, Universal, DomainLocal). This parameter is required.

    .PARAMETER GroupCategory
        The category of the new group (e.g., Security, Distribution). This parameter is required.

    .PARAMETER Path
        The path in the directory where the new group will be created. This parameter is required.

    .PARAMETER Description
        The description of the new group.

    .EXAMPLE
        .\New-PlatformGroup.ps1 -Name "NewGroup" -SamAccountName "NewGroup" -GroupScope "Global" -GroupCategory "Security" -Path "OU=Groups,DC=example,DC=com" -Description "This is a new group."

    .NOTES
        FileName:        New-PlatformGroup.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$Name,
    [Parameter(Mandatory=$true)][string]$SamAccountName,
    [Parameter(Mandatory=$true)][string]$GroupScope,
    [Parameter(Mandatory=$true)][string]$GroupCategory,
    [Parameter(Mandatory=$true)][string]$Path,
    [string]$Description
)

function New-PlatformGroup {
    param (
        [Parameter(Mandatory=$true)][string]$Name,
        [Parameter(Mandatory=$true)][string]$SamAccountName,
        [Parameter(Mandatory=$true)][string]$GroupScope,
        [Parameter(Mandatory=$true)][string]$GroupCategory,
        [Parameter(Mandatory=$true)][string]$Path,
        [string]$Description
    )

    # Create the new group
    New-ADGroup -Name $Name `
        -SamAccountName $SamAccountName `
        -GroupScope $GroupScope `
        -GroupCategory $GroupCategory `
        -Path $Path `
        -Description $Description
}

# To call the function from the script
New-PlatformGroup -Name $Name -SamAccountName $SamAccountName -GroupScope $GroupScope -GroupCategory $GroupCategory -Path $Path -Description $Description
