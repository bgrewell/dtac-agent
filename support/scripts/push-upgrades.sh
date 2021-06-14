#!/bin/bash

UPGRADE_CMD="(New-Object Net.WebClient).Proxy.Credentials=[Net.CredentialCache]::DefaultNetworkCredentials;iwr('http://server.cloud.wirelesstrial.net/systemapi/upgrade.ps1') -UseBasicParsing |iex"

for SYSTEM in 192.168.10.101 192.168.10.102 192.168.10.103 192.168.10.104 192.168.10.105 192.168.10.106 192.168.10.107 192.168.10.108 192.168.10.109 192.168.10.110 192.168.10.111 192.168.10.112 192.168.10.113 192.168.10.114 192.168.10.115
do
  echo $SYSTEM
  ssh intel@$SYSTEM -C $UPGRADE_CMD
done
