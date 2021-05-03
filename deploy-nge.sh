#!/usr/bin/env bash

echo "building latest version"
make build

systems=( 10.108.1.21 10.108.2.61 10.108.1.11 )

for system in "${systems[@]}"
do
  echo "deploying to $system"
  scp bin/system-apid intel@$system:/home/intel/.
  scp support/service/system-apid.service intel@$system:/home/intel/.
  scp support/config/config.yaml intel@$system:/home/intel/.
  ssh intel@$system -C 'sudo systemctl stop system-apid || true'
  ssh intel@$system -C 'sudo mkdir -p /opt/system-api/bin || true'
  ssh intel@$system -C 'sudo mkdir -p /etc/system-api || true'
  ssh intel@$system -C 'sudo mv ~/system-apid /opt/system-api/bin/.'
  ssh intel@$system -C 'sudo mv ~/system-apid.service /lib/systemd/system/.'
  ssh intel@$system -C 'sudo mv ~/config.yaml /etc/system-api/config.yaml'
  ssh intel@$system -C 'sudo systemctl daemon-reload'
  ssh intel@$system -C 'sudo systemctl start system-apid'
  echo "finished deploying to $system"
done