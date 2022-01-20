#!/bin/bash

UPGRADE_CMD="(New-Object Net.WebClient).Proxy.Credentials=[Net.CredentialCache]::DefaultNetworkCredentials;iwr('http://server.cloud.wirelesstrial.net/systemapi/upgrade.ps1') -UseBasicParsing |iex"

for SYSTEM in 192.168.10.101 192.168.10.102 192.168.10.103 192.168.10.104 192.168.10.105 192.168.10.106 192.168.10.107 192.168.10.108 192.168.10.109 192.168.10.110 192.168.10.111 192.168.10.112 192.168.10.113 192.168.10.114 192.168.10.115
do
  echo $SYSTEM
  ssh intel@$SYSTEM -C $UPGRADE_CMD
done

SERVERS=( 10.108.1.21 10.108.2.61 10.108.1.11 )

for server in "${SERVERS[@]}"
do
  echo "deploying to $server"
  scp bin/system-apid intel@$server:/home/intel/.
  scp support/service/system-apid.service intel@$server:/home/intel/.
  scp support/config/config.yaml intel@$server:/home/intel/.
  ssh intel@$server -C 'sudo systemctl stop system-apid || true'
  ssh intel@$server -C 'sudo mkdir -p /opt/system-api/bin || true'
  ssh intel@$server -C 'sudo mkdir -p /etc/system-api || true'
  ssh intel@$server -C 'sudo mv ~/system-apid /opt/system-api/bin/.'
  ssh intel@$server -C 'sudo mv ~/system-apid.service /lib/systemd/system/.'
  ssh intel@$server -C 'sudo mv ~/config.yaml /etc/system-api/config.yaml'
  ssh intel@$server -C 'sudo systemctl daemon-reload'
  ssh intel@$server -C 'sudo systemctl start system-apid'
  echo "finished deploying to $server"
done
