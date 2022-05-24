Write-Output "[+] === Creating directory ==="
mkdir "C:\Program Files\Intel\DTAC-Agent"
Write-Output "[+] === Downloading dtac-agentd.exe ==="
Invoke-WebRequest -Uri http://software.labs.intel.com/dtac/agent/windows/dtac-agent.exe -OutFile "C:\Program Files\Intel\DTAC-Agent\dtac-agent.exe"
Write-Output "[+] === Downloading DTAC Agent configuration ==="
Invoke-WebRequest -Uri http://server.cloud.wirelesstrial.net/systemapi/config.yaml -OutFile "C:\Program Files\Intel\System-Api\config.yaml"
Write-Output "[+] === Creating firewall rule ==="
New-NetFirewallRule -DisplayName "DTAC-Agent" -Profile Any -Direction Inbound -Action Allow -Protocol TCP -LocalPort 8080
Write-Output "[+] === Enabling service ==="
cd "C:\Program Files\Intel\DTAC-Agent"
.\dtac-agentd.exe --service install
Start-Service -Name dtac-agent.service
Write-Output "[!] === DONE ==="
