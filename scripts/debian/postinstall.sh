#!/usr/bin/env bash

# Ensure Python 3 and pip are installed
if ! command -v python3 &> /dev/null
then
    echo "Python 3 is not installed. Installing Python 3..."
    apt-get update
    apt-get install -y python3
fi

if ! command -v pip3 &> /dev/null
then
    echo "pip for Python 3 is not installed. Installing pip3..."
    apt-get update
    apt-get install -y python3-pip
fi

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

# Restart the service
systemctl restart dtac-agentd.service