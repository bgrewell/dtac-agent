<#
    .SYNOPSIS
        Get-PlatformGroup - Retrieve group information from Active Directory and export to CSV.

    .DESCRIPTION
        This script retrieves group information from Active Directory based on either a specific SamAccountName or a filter.
        It exports the retrieved information to a CSV format.

    .PARAMETER SamAccountName
        The SamAccountName of the group to retrieve.

    .PARAMETER Filter
        The filter to apply when retrieving groups. Defaults to '*', which retrieves all groups.

    .EXAMPLE
        .\Get-PlatformGroup.ps1 -SamAccountName GroupName

    .EXAMPLE
        .\Get-PlatformGroup.ps1 -Filter "GroupCategory -eq 'Security'"

    .NOTES
        FileName:        Get-PlatformGroup.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [string]$SamAccountName,
    [string]$Filter = "*"
)

function Get-PlatformGroup {
    param (
        [string]$SamAccountName,
        [string]$Filter = "*"
    )

    if ($SamAccountName) {
        # Get a single group and export to CSV
        Get-ADGroup -Identity $SamAccountName -Properties Name, SamAccountName, DistinguishedName, GroupCategory, GroupScope, SID |
        Select-Object Name, SamAccountName, DistinguishedName, GroupCategory, GroupScope, SID |
        ConvertTo-Csv -NoTypeInformation
    } else {
        # Get all groups or groups matching the filter and export to CSV
        Get-ADGroup -Filter $Filter -Properties Name, SamAccountName, DistinguishedName, GroupCategory, GroupScope, SID |
        Select-Object Name, SamAccountName, DistinguishedName, GroupCategory, GroupScope, SID |
        ConvertTo-Csv -NoTypeInformation
    }
}

# To call the function from the script
if ($SamAccountName) {
    Get-PlatformGroup -SamAccountName $SamAccountName
} else {
    Get-PlatformGroup -Filter $Filter
}
