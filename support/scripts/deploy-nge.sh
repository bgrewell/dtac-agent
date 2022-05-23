#!/usr/bin/env bash

echo "building latest version"
make build

systems=( 10.108.1.21 10.108.2.61 10.108.1.11 )

for system in "${systems[@]}"
do
  echo "deploying to $system"
  scp bin/system-agentd intel@$system:/home/intel/.
  scp support/service/system-agentd.service intel@$system:/home/intel/.
  scp support/config/config.yaml intel@$system:/home/intel/.
  ssh intel@$system -C 'sudo systemctl stop system-agentd || true'
  ssh intel@$system -C 'sudo mkdir -p /opt/system-agent/bin || true'
  ssh intel@$system -C 'sudo mkdir -p /etc/system-agent || true'
  ssh intel@$system -C 'sudo mv ~/system-agentd /opt/system-agent/bin/.'
  ssh intel@$system -C 'sudo mv ~/system-agentd.service /lib/systemd/system/.'
  ssh intel@$system -C 'sudo mv ~/config.yaml /etc/system-agent/config.yaml'
  ssh intel@$system -C 'sudo systemctl daemon-reload'
  ssh intel@$system -C 'sudo systemctl start system-agentd'
  echo "finished deploying to $system"
done