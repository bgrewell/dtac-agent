<#
    .SYNOPSIS
        Get-PlatformUser - Retrieve user information from Active Directory and export to CSV.

    .DESCRIPTION
        This script retrieves user information from Active Directory based on either a specific SamAccountName or a filter.
        It exports the retrieved information to a CSV format.

    .PARAMETER SamAccountName
        The SamAccountName of the user to retrieve.

    .PARAMETER Filter
        The filter to apply when retrieving users. Defaults to '*', which retrieves all users.

    .EXAMPLE
        .\Get-PlatformUser.ps1 -SamAccountName jdoe

    .EXAMPLE
        .\Get-PlatformUser.ps1 -Filter "Department -eq 'Sales'"

    .NOTES
        FileName:        Get-PlatformUser.ps1
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

function Get-PlatformUser {
    param (
        [string]$SamAccountName,
        [string]$Filter = "*"
    )

    if ($SamAccountName) {
        # Get a single user and export to CSV
        Get-ADUser -Identity $SamAccountName -Properties DisplayName, SamAccountName, UserPrincipalName, EmailAddress, Enabled, SSHPublicKey |
        Select-Object DisplayName, SamAccountName, UserPrincipalName, EmailAddress, Enabled, @{Name='SSHPublicKey';Expression={$_.SSHPublicKey -join ','}} |
        ConvertTo-Csv -NoTypeInformation
    } else {
        # Get all users or users matching the filter and export to CSV
        Get-ADUser -Filter $Filter -Properties DisplayName, SamAccountName, UserPrincipalName, EmailAddress, Enabled, SSHPublicKey |
        Select-Object DisplayName, SamAccountName, UserPrincipalName, EmailAddress, Enabled, @{Name='SSHPublicKey';Expression={$_.SSHPublicKey -join ','}} |
        ConvertTo-Csv -NoTypeInformation
    }
}

# To call the function from the script
if ($SamAccountName) {
    Get-PlatformUser -SamAccountName $SamAccountName
} else {
    Get-PlatformUser -Filter $Filter
}
