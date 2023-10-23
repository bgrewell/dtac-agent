#!/usr/bin/env bash

# Ensure that it doesn't already exist
multipass delete dtac-test
multipass purge

# Create VM
multipass launch -n dtac-test

# Copy the files
multipass transfer ../../dist/*.deb dtac-test:/home/ubuntu/dtac-install.deb

# Execute
multipass exec dtac-test -- sudo dpkg -i /home/ubuntu/dtac-install.deb