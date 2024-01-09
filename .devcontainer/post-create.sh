#!/usr/bin/env bash

# Setup mage
git clone https://github.com/magefile/mage /tmp/mage
cd /tmp/mage
go run bootstrap.go

# Tell git to ignore ownership issues
git config --global --add safe.directory /workspaces/dtac-agent

# Create bin and plugins directory
mkdir -p /opt/dtac/{plugins,bin}

# Install tools
apt update
apt install -y curl vim wget tree