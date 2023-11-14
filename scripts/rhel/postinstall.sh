#!/usr/bin/env bash

# Install dtac-tools globally using pip3
pip3 install dtac-tools

# install yq
VERSION=v4.2.0
BINARY=yq_linux_amd64
wget https://github.com/mikefarah/yq/releases/download/${VERSION}/${BINARY} -O /usr/bin/yq
chmod +x /usr/bin/yq

# Reload the daemon
systemctl daemon-reload

# Enable the service
systemctl enable dtac-agentd.service

# Start the service
systemctl start dtac-agentd.service

# Generates a random password
password=$(openssl rand -base64 32)

# Updates the password in the YAML file
yq eval -i '.authn.pass = "'"$password"'"' /etc/dtac/config.yaml

# Generate a link between the dtac config utility and /usr/bin
ln -s /opt/dtac/bin/dtac /usr/bin/dtac

# Move the plugins
mkdir -p /opt/dtac/plugins
mv /opt/dtac/bin/*.plugin /opt/dtac/plugins/.

# Restart the service
systemctl restart dtac-agentd.service