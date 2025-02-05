<#
    .SYNOPSIS
        Get-PlatformSudoRole - Retrieve sudo role information from Active Directory and export to CSV.

    .DESCRIPTION
        This script retrieves sudo role information from Active Directory based on either a specific RoleName or a filter.
        It exports the retrieved information to a CSV format.

    .PARAMETER RoleName
        The RoleName of the sudo role to retrieve.

    .PARAMETER Path
        The path in the directory where the sudo roles are located. This parameter is required.

    .PARAMETER Filter
        The filter to apply when retrieving sudo roles. Defaults to '*', which retrieves all roles.

    .EXAMPLE
        .\Get-PlatformSudoRole.ps1 -RoleName SudoRoleName -Path "OU=Sudoers,DC=example,DC=com"

    .EXAMPLE
        .\Get-PlatformSudoRole.ps1 -Filter "SudoUser -eq 'admin'" -Path "OU=Sudoers,DC=example,DC=com"

    .NOTES
        FileName:        Get-PlatformSudoRole.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [string]$RoleName,
    [Parameter(Mandatory=$true)][string]$Path,
    [string]$Filter = "*"
)

function Get-PlatformSudoRole {
    param (
        [string]$RoleName,
        [string]$Path,
        [string]$Filter = "*"
    )

    if ($RoleName) {
        # Get a single sudo role and export to CSV
        Get-ADObject -Filter "Name -eq '$RoleName'" -SearchBase $Path -Properties Name, Path, SudoHost, SudoUser, SudoCommand |
        Select-Object Name, Path, @{Name='SudoHost';Expression={$_.SudoHost -join ','}}, @{Name='SudoUser';Expression={$_.SudoUser -join ','}}, @{Name='SudoCommand';Expression={$_.SudoCommand -join ','}} |
        ConvertTo-Csv -NoTypeInformation
    } else {
        # Get all sudo roles or roles matching the filter and export to CSV
        Get-ADObject -Filter $Filter -SearchBase $Path -Properties Name, Path, SudoHost, SudoUser, SudoCommand |
        Select-Object Name, Path, @{Name='SudoHost';Expression={$_.SudoHost -join ','}}, @{Name='SudoUser';Expression={$_.SudoUser -join ','}}, @{Name='SudoCommand';Expression={$_.SudoCommand -join ','}} |
        ConvertTo-Csv -NoTypeInformation
    }
}

# To call the function from the script
if ($RoleName) {
    Get-PlatformSudoRole -RoleName $RoleName -Path $Path
} else {
    Get-PlatformSudoRole -Filter $Filter -Path $Path
}
