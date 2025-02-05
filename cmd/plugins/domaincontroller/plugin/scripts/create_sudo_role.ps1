<#
    .SYNOPSIS
        New-PlatformSudoRole - Create a new sudo role object in Active Directory.

    .DESCRIPTION
        This script creates a new sudo role object in Active Directory with the specified parameters.

    .PARAMETER RoleName
        The name of the new sudo role. This parameter is required.

    .PARAMETER Path
        The path in the directory where the new sudo role will be created. This parameter is required.

    .PARAMETER SudoHost
        The host(s) for which the sudo role applies. This parameter is required.

    .PARAMETER SudoUser
        The user(s) for which the sudo role applies. This parameter is required.

    .PARAMETER SudoCommand
        The command(s) for which the sudo role applies. Defaults to 'ALL'.

    .EXAMPLE
        .\New-PlatformSudoRole.ps1 -RoleName "AdminRole" -Path "OU=Sudoers,DC=example,DC=com" -SudoHost "ALL" -SudoUser "admin" -SudoCommand "/usr/bin/passwd"

    .NOTES
        FileName:        New-PlatformSudoRole.ps1
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
    [Parameter(Mandatory=$true)][string]$SudoHost,
    [Parameter(Mandatory=$true)][string]$SudoUser,
    [string]$SudoCommand = 'ALL'
)

function New-PlatformSudoRole {
    param (
        [Parameter(Mandatory=$true)][string]$RoleName,
        [Parameter(Mandatory=$true)][string]$Path,
        [Parameter(Mandatory=$true)][string]$SudoHost,
        [Parameter(Mandatory=$true)][string]$SudoUser,
        [string]$SudoCommand = 'ALL'
    )

    # Create the new sudo role object
    New-ADObject -Name $RoleName `
        -Path $Path `
        -Type 'sudoRole' `
        -OtherAttributes @{
            sudoCommand = $SudoCommand;
            sudoHost    = $SudoHost;
            sudoUser    = $SudoUser
        }
}

# To call the function from the script
New-PlatformSudoRole -RoleName $RoleName -Path $Path -SudoHost $SudoHost -SudoUser $SudoUser -SudoCommand $SudoCommand
