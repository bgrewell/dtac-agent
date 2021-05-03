@echo off

echo "building for linux"
set GOARCH=amd64
set GOOS=linux
go build -o bin/system-apid main.go

set systems=10.108.1.11 10.108.1.21 10.108.2.61

(for %%a in (%systems%) do (
    echo deploying System-API to %%a
    scp bin\system-apid intel@%%a:/home/intel/.
    scp support\service\system-apid.service intel@%%a:/home/intel/.
    ssh intel@%%a -C 'sudo systemctl stop system-apid || true'
    ssh intel@%%a -C 'sudo mkdir -p /opt/system-api/bin || true'
    ssh intel@%%a -C 'sudo mv ~/system-apid /opt/system-api/bin/.'
    ssh intel@%%a -C 'sudo mv ~/system-apid.service /lib/systemd/system/.'
    ssh intel@%%a -C 'sudo systemctl daemon-reload'
    ssh intel@%%a -C 'sudo systemctl start system-apid'
    echo finished updating %%a
))