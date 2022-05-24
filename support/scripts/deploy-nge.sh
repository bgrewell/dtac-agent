#!/usr/bin/env bash

echo "building latest version"
make build

systems=( 10.108.1.21 10.108.2.61 10.108.1.11 )

for system in "${systems[@]}"
do
  echo "deploying to $system"
  scp bin/dtac-agentd intel@$system:/home/intel/.
  scp support/service/dtac-agentd.service intel@$system:/home/intel/.
  scp support/config/config.yaml intel@$system:/home/intel/.
  ssh intel@$system -C 'sudo systemctl stop dtac-agentd || true'
  ssh intel@$system -C 'sudo mkdir -p /opt/dtac-agent/bin || true'
  ssh intel@$system -C 'sudo mkdir -p /etc/dtac-agent || true'
  ssh intel@$system -C 'sudo mv ~/dtac-agentd /opt/dtac-agent/bin/.'
  ssh intel@$system -C 'sudo mv ~/dtac-agentd.service /lib/systemd/system/.'
  ssh intel@$system -C 'sudo mv ~/config.yaml /etc/dtac-agent/config.yaml'
  ssh intel@$system -C 'sudo systemctl daemon-reload'
  ssh intel@$system -C 'sudo systemctl start dtac-agentd'
  echo "finished deploying to $system"
done