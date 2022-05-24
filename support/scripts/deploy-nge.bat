@echo off

echo "building for linux"
set GOARCH=amd64
set GOOS=linux
go build -o bin/dtac-agentd main.go

set systems=10.108.1.11 10.108.1.21 10.108.2.61

(for %%a in (%systems%) do (
    echo deploying dtac-agent to %%a
    scp bin\dtac-agentd intel@%%a:/home/intel/.
    scp support\service\dtac-agentd.service intel@%%a:/home/intel/.
    ssh intel@%%a -C 'sudo systemctl stop dtac-agentd || true'
    ssh intel@%%a -C 'sudo mkdir -p /opt/dtac-agent/bin || true'
    ssh intel@%%a -C 'sudo mv ~/dtac-agentd /opt/dtac-agent/bin/.'
    ssh intel@%%a -C 'sudo mv ~/dtac-agentd.service /lib/systemd/system/.'
    ssh intel@%%a -C 'sudo systemctl daemon-reload'
    ssh intel@%%a -C 'sudo systemctl start dtac-agentd'
    echo finished updating %%a
))