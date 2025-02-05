<#
    .SYNOPSIS
        New-PlatformUser - Create a new user in Active Directory.

    .DESCRIPTION
        This script creates a new user in Active Directory with the specified parameters.

    .PARAMETER Name
        The name of the new user. This parameter is required.

    .PARAMETER GivenName
        The given name of the new user. This parameter is required.

    .PARAMETER Surname
        The surname of the new user. This parameter is required.

    .PARAMETER UserPrincipalName
        The User Principal Name (UPN) of the new user. This parameter is required.

    .PARAMETER SamAccountName
        The SamAccountName of the new user. This parameter is required.

    .PARAMETER Path
        The path in the directory where the new user will be created. This parameter is required.

    .PARAMETER Password
        The password for the new user. This parameter is required.

    .PARAMETER Description
        The description of the new user.

    .PARAMETER EmailAddress
        The email address of the new user.

    .PARAMETER Department
        The department of the new user.

    .PARAMETER Title
        The title of the new user.

    .EXAMPLE
        .\New-PlatformUser.ps1 -Name "John Doe" -GivenName "John" -Surname "Doe" -UserPrincipalName "johndoe@example.com" -SamAccountName "johndoe" -Path "OU=Users,DC=example,DC=com" -Password "P@ssw0rd" -Description "New user account" -EmailAddress "johndoe@example.com" -Department "IT" -Title "Developer"

    .NOTES
        FileName:        New-PlatformUser.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$Name,
    [Parameter(Mandatory=$true)][string]$GivenName,
    [Parameter(Mandatory=$true)][string]$Surname,
    [Parameter(Mandatory=$true)][string]$UserPrincipalName,
    [Parameter(Mandatory=$true)][string]$SamAccountName,
    [Parameter(Mandatory=$true)][string]$Path,
    [Parameter(Mandatory=$true)][string]$Password,
    [string]$Initials,
    [string]$Description,
    [string]$EmailAddress,
    [string]$Department,
    [string]$Title
)

function New-PlatformUser {
    param (
        [Parameter(Mandatory=$true)][string]$Name,
        [Parameter(Mandatory=$true)][string]$GivenName,
        [Parameter(Mandatory=$true)][string]$Surname,
        [Parameter(Mandatory=$true)][string]$UserPrincipalName,
        [Parameter(Mandatory=$true)][string]$SamAccountName,
        [Parameter(Mandatory=$true)][string]$Path,
        [Parameter(Mandatory=$true)][string]$Password,
        [string]$Initials,
        [string]$Description,
        [string]$EmailAddress,
        [string]$Department,
        [string]$Title
    )

    # Create the new user
    New-ADUser -Name $Name `
        -GivenName $GivenName `
        -Surname $Surname `
        -Initials $Initials `
        -UserPrincipalName $UserPrincipalName `
        -SamAccountName $SamAccountName `
        -Path $Path `
        -AccountPassword (ConvertTo-SecureString $Password -AsPlainText -Force) `
        -Enabled $true `
        -Description $Description `
        -EmailAddress $EmailAddress `
        -Department $Department `
        -Title $Title
}

# To call the function from the script
New-PlatformUser -Name $Name -GivenName $GivenName -Surname $Surname -UserPrincipalName $UserPrincipalName -SamAccountName $SamAccountName -Path $Path -Password $Password -Description $Description -EmailAddress $EmailAddress -Department $Department -Title $Title -Initials $Initials
