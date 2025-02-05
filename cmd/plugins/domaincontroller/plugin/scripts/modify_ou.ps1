<#
    .SYNOPSIS
        Set-PlatformOrganizationalUnit - Modify an organizational unit (OU) in Active Directory.

    .DESCRIPTION
        This script modifies an organizational unit (OU) in Active Directory based on the specified parameters.

    .PARAMETER Identity
        The distinguished name (DN) or GUID of the OU to be modified. This parameter is required.

    .PARAMETER Name
        The new name of the OU.

    .PARAMETER Description
        The new description of the OU.

    .PARAMETER ProtectedFromAccidentalDeletion
        Specifies whether the OU is protected from accidental deletion. This parameter is required.

    .PARAMETER ManagedBy
        The distinguished name (DN) of the user or group that manages this OU.

    .PARAMETER StreetAddress
        The street address of the OU.

    .PARAMETER City
        The city of the OU.

    .PARAMETER State
        The state or province of the OU.

    .PARAMETER Country
        The country of the OU.

    .PARAMETER PostalCode
        The postal code of the OU.

    .PARAMETER PhoneNumber
        The phone number of the OU.

    .PARAMETER EmailAddress
        The email address of the OU.

    .PARAMETER Fax
        The fax number of the OU.

    .EXAMPLE
        .\Set-PlatformOrganizationalUnit.ps1 -Identity "OU=Departments,DC=example,DC=com" -Name "NewOU" -Description "Updated Description" -ProtectedFromAccidentalDeletion $true -ManagedBy "CN=Manager,OU=Users,DC=example,DC=com" -StreetAddress "123 Main St" -City "Anytown" -State "AnyState" -Country "AnyCountry" -PostalCode "12345" -PhoneNumber "555-5555" -EmailAddress "ou@example.com" -Fax "555-5556"

    .NOTES
        FileName:        Set-PlatformOrganizationalUnit.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$Identity,
    [string]$Name,
    [string]$Description,
    [bool]$ProtectedFromAccidentalDeletion,
    [string]$ManagedBy,
    [string]$StreetAddress,
    [string]$City,
    [string]$State,
    [string]$Country,
    [string]$PostalCode,
    [string]$PhoneNumber,
    [string]$EmailAddress,
    [string]$Fax
)

function Set-PlatformOrganizationalUnit {
    param (
        [Parameter(Mandatory=$true)][string]$Identity,
        [string]$Name,
        [string]$Description,
        [bool]$ProtectedFromAccidentalDeletion,
        [string]$ManagedBy,
        [string]$StreetAddress,
        [string]$City,
        [string]$State,
        [string]$Country,
        [string]$PostalCode,
        [string]$PhoneNumber,
        [string]$EmailAddress,
        [string]$Fax
    )

    # Modify the organizational unit
    Set-ADOrganizationalUnit -Identity $Identity `
        -Name $Name `
        -Description $Description `
        -ManagedBy $ManagedBy `
        -StreetAddress $StreetAddress `
        -City $City `
        -State $State `
        -Country $Country `
        -PostalCode $PostalCode `
        -PhoneNumber $PhoneNumber `
        -EmailAddress $EmailAddress `
        -Fax $Fax `
        -ProtectedFromAccidentalDeletion $ProtectedFromAccidentalDeletion
}

# To call the function from the script
Set-PlatformOrganizationalUnit -Identity $Identity -Name $Name -Description $Description -ProtectedFromAccidentalDeletion $ProtectedFromAccidentalDeletion -ManagedBy $ManagedBy -StreetAddress $StreetAddress -City $City -State $State -Country $Country -PostalCode $PostalCode -PhoneNumber $PhoneNumber -EmailAddress $EmailAddress -Fax $Fax
