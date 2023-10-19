#!/usr/bin/env bash

# Reload the daemon
sudo systemctl daemon-reload

# Enable the service
sudo systemctl enable dtac-agentd.service

# Start the service
sudo systemctl start dtac-agentd.service