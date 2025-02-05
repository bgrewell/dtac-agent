<#
    .SYNOPSIS
        Rename-PlatformOrganizationalUnit - Rename and move an organizational unit (OU) in Active Directory.

    .DESCRIPTION
        This script renames and moves an organizational unit (OU) in Active Directory based on the specified parameters.

    .PARAMETER Identity
        The distinguished name (DN) or GUID of the OU to be modified. This parameter is required.

    .PARAMETER NewName
        The new name of the OU. This parameter is required.

    .PARAMETER NewPath
        The new path in the directory where the OU will be moved. This parameter is required.

    .EXAMPLE
        .\Rename-PlatformOrganizationalUnit.ps1 -Identity "OU=OldOU,DC=example,DC=com" -NewName "NewOU" -NewPath "OU=Departments,DC=example,DC=com"

    .NOTES
        FileName:        Rename-PlatformOrganizationalUnit.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$Identity,
    [Parameter(Mandatory=$true)][string]$NewName,
    [Parameter(Mandatory=$true)][string]$NewPath
)

function Rename-PlatformOrganizationalUnit {
    param (
        [Parameter(Mandatory=$true)][string]$Identity,
        [Parameter(Mandatory=$true)][string]$NewName,
        [Parameter(Mandatory=$true)][string]$NewPath
    )

    # Construct the new distinguished name (DN)
    $parentPath = (Get-ADOrganizationalUnit -Identity $NewPath).DistinguishedName
    $newDn = "OU=$NewName,$parentPath"

    Rename-ADObject -Identity $Identity -NewName $NewName -TargetPath $parentPath
}

# To call the function from the script
Rename-PlatformOrganizationalUnit -Identity $Identity -NewName $NewName -NewPath $NewPath