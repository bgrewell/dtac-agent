<#
    .SYNOPSIS
        Get-PlatformOU - Retrieve organizational unit (OU) information from Active Directory and export to CSV.

    .DESCRIPTION
        This script retrieves OU information from Active Directory based on either a specific DistinguishedName or a filter.
        It exports the retrieved information to a CSV format.

    .PARAMETER DistinguishedName
        The DistinguishedName of the OU to retrieve.

    .PARAMETER Filter
        The filter to apply when retrieving OUs. Defaults to '*', which retrieves all OUs.

    .EXAMPLE
        .\Get-PlatformOU.ps1 -DistinguishedName "OU=ExampleOU,DC=example,DC=com"

    .EXAMPLE
        .\Get-PlatformOU.ps1 -Filter "Description -like '*Department*'"

    .NOTES
        FileName:        Get-PlatformOU.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [string]$DistinguishedName,
    [string]$Filter = "*"
)

function Get-PlatformOU {
    param (
        [string]$DistinguishedName,
        [string]$Filter = "*"
    )

    if ($DistinguishedName) {
        # Get a single OU and export to CSV
        Get-ADOrganizationalUnit -Identity $DistinguishedName -Properties Name, DistinguishedName, Description, ObjectGUID |
        Select-Object Name, DistinguishedName, Description, ObjectGUID |
        ConvertTo-Csv -NoTypeInformation
    } else {
        # Get all OUs or OUs matching the filter and export to CSV
        Get-ADOrganizationalUnit -Filter $Filter -Properties Name, DistinguishedName, Description, ObjectGUID |
        Select-Object Name, DistinguishedName, Description, ObjectGUID |
        ConvertTo-Csv -NoTypeInformation
    }
}

# To call the function from the script
if ($DistinguishedName) {
    Get-PlatformOU -DistinguishedName $DistinguishedName
} else {
    Get-PlatformOU -Filter $Filter
}
