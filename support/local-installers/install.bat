:: This script is designed to install dtac-agent service on Windows systems
:: It should work out of the box for most systems, but modify as needed if
:: your use case requires customization.
:: Reach out to Benjamin Grewell <benjamin.grewell@intel.com> with any questions
set "params=%*"
cd /d "%~dp0" && ( if exist "%temp%\getadmin.vbs" del "%temp%\getadmin.vbs" ) && fsutil dirty query %systemdrive% 1>nul 2>nul || (  echo Set UAC = CreateObject^("Shell.Application"^) : UAC.ShellExecute "cmd.exe", "/k cd ""%~sdp0"" && %~s0 %params%", "", "runas", 1 >> "%temp%\getadmin.vbs" && "%temp%\getadmin.vbs" && exit /B )

:: Create directories
mkdir "C:\Program Files\Intel\dtac-agent\logs"