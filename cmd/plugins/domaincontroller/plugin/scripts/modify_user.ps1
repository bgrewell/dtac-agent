<#
    .SYNOPSIS
        Set-PlatformUser - Modify a user in Active Directory.

    .DESCRIPTION
        This script modifies a user in Active Directory based on the specified parameters.

    .PARAMETER SamAccountName
        The SamAccountName of the user to be modified. This parameter is required.

    .PARAMETER GivenName
        The new given name of the user.

    .PARAMETER Surname
        The new surname of the user.

    .PARAMETER UserPrincipalName
        The new User Principal Name (UPN) of the user.

    .PARAMETER Description
        The new description of the user.

    .PARAMETER EmailAddress
        The new email address of the user.

    .PARAMETER Department
        The new department of the user.

    .PARAMETER Title
        The new title of the user.

    .EXAMPLE
        .\Set-PlatformUser.ps1 -SamAccountName "jdoe" -GivenName "John" -Surname "Doe" -UserPrincipalName "johndoe@example.com" -Description "Updated Description" -EmailAddress "johndoe@example.com" -Department "IT" -Title "Senior Developer"

    .NOTES
        FileName:        Set-PlatformUser.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$SamAccountName,
    [string]$GivenName,
    [string]$Surname,
    [string]$UserPrincipalName,
    [string]$Description,
    [string]$EmailAddress,
    [string]$Department,
    [string]$Title
)

function Set-PlatformUser {
    param (
        [Parameter(Mandatory=$true)][string]$SamAccountName,
        [string]$GivenName,
        [string]$Surname,
        [string]$UserPrincipalName,
        [string]$Description,
        [string]$EmailAddress,
        [string]$Department,
        [string]$Title
    )

    # Construct the hash table for the properties to be updated
    $properties = @{}

    if ($GivenName) { $properties["GivenName"] = $GivenName }
    if ($Surname) { $properties["Surname"] = $Surname }
    if ($UserPrincipalName) { $properties["UserPrincipalName"] = $UserPrincipalName }
    if ($Description) { $properties["Description"] = $Description }
    if ($EmailAddress) { $properties["EmailAddress"] = $EmailAddress }
    if ($Department) { $properties["Department"] = $Department }
    if ($Title) { $properties["Title"] = $Title }

    # Modify the user
    if ($properties.Count -gt 0) {
        Set-ADUser -Identity $SamAccountName @properties
    } else {
        Write-Host "No properties to update for $SamAccountName"
    }
}

# To call the function from the script
Set-PlatformUser -SamAccountName $SamAccountName -GivenName $GivenName -Surname $Surname -UserPrincipalName $UserPrincipalName -Description $Description -EmailAddress $EmailAddress -Department $Department -Title $Title
