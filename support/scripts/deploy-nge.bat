@echo off

echo "building for linux"
set GOARCH=amd64
set GOOS=linux
go build -o bin/system-agentd main.go

set systems=10.108.1.11 10.108.1.21 10.108.2.61

(for %%a in (%systems%) do (
    echo deploying system-agent to %%a
    scp bin\system-agentd intel@%%a:/home/intel/.
    scp support\service\system-agentd.service intel@%%a:/home/intel/.
    ssh intel@%%a -C 'sudo systemctl stop system-agentd || true'
    ssh intel@%%a -C 'sudo mkdir -p /opt/system-agent/bin || true'
    ssh intel@%%a -C 'sudo mv ~/system-apid /opt/system-api/bin/.'
    ssh intel@%%a -C 'sudo mv ~/system-apid.service /lib/systemd/system/.'
    ssh intel@%%a -C 'sudo systemctl daemon-reload'
    ssh intel@%%a -C 'sudo systemctl start system-apid'
    echo finished updating %%a
))