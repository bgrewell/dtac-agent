<#
    .SYNOPSIS
        Set-PlatformSudoRole - Modify a sudo role object in Active Directory.

    .DESCRIPTION
        This script modifies a sudo role object in Active Directory based on the specified parameters.

    .PARAMETER RoleName
        The name of the sudo role to be modified. This parameter is required.

    .PARAMETER Path
        The path in the directory where the sudo role is located. This parameter is required.

    .PARAMETER SudoHost
        The new host(s) for which the sudo role applies.

    .PARAMETER SudoUser
        The new user(s) for which the sudo role applies.

    .PARAMETER SudoCommand
        The new command(s) for which the sudo role applies.

    .EXAMPLE
        .\Set-PlatformSudoRole.ps1 -RoleName "AdminRole" -Path "OU=Sudoers,DC=example,DC=com" -SudoHost "ALL" -SudoUser "admin" -SudoCommand "/usr/bin/passwd"

    .NOTES
        FileName:        Set-PlatformSudoRole.ps1
        Author:          Benjamin Grewell
        Email:           benjamin.grewell@intel.com
        Last Updated:    28-May-2024
        License:         GPL (GNU General Public License)

    .LINK
        https://www.gnu.org/licenses/gpl-3.0.en.html
#>

param (
    [Parameter(Mandatory=$true)][string]$RoleName,
    [Parameter(Mandatory=$true)][string]$Path,
    [string]$SudoHost,
    [string]$SudoUser,
    [string]$SudoCommand
)

function Set-PlatformSudoRole {
    param (
        [Parameter(Mandatory=$true)][string]$RoleName,
        [Parameter(Mandatory=$true)][string]$Path,
        [string]$SudoHost,
        [string]$SudoUser,
        [string]$SudoCommand
    )

    # Get the distinguished name of the sudo role object
    $sudoRoleDN = (Get-ADObject -Filter "Name -eq '$RoleName'" -SearchBase $Path).DistinguishedName

    # Modify the sudo role object
    Set-ADObject -Identity $sudoRoleDN `
        -Replace @{
            sudoCommand = $SudoCommand;
            sudoHost    = $SudoHost;
            sudoUser    = $SudoUser
        }
}

# To call the function from the script
Set-PlatformSudoRole -RoleName $RoleName -Path $Path -SudoHost $SudoHost -SudoUser $SudoUser -SudoCommand $SudoCommand
