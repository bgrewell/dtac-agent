#!/usr/bin/env bash

# Set defaults if variables are not set
DEPLOY_SERVER="${DEPLOY_SERVER:-localhost}"
DEPLOY_USER="${DEPLOY_USER:-$USER}"

echo "Attempting ssh connection $DEPLOY_USER@$DEPLOY_SERVER"
scp bin/dtac-agentd.exe $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/windows/.
ssh $DEPLOY_USER@$DEPLOY_SERVER 'md5sum /var/www/html/dtac/agent/windows/dtac-agentd.exe > /var/www/html/dtac/agent/windows/md5'
ssh $DEPLOY_USER@$DEPLOY_SERVER 'sha256sum /var/www/html/dtac/agent/windows/dtac-agentd.exe > /var/www/html/dtac/agent/windows/sha256'
scp bin/dtac-agentd $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/linux/.
ssh $DEPLOY_USER@$DEPLOY_SERVER 'md5sum /var/www/html/dtac/agent/linux/dtac-agentd > /var/www/html/dtac/agent/linux/md5'
ssh $DEPLOY_USER@$DEPLOY_SERVER 'sha256sum /var/www/html/dtac/agent/linux/dtac-agentd > /var/www/html/dtac/agent/linux/sha256'
scp bin/dtac-agentd.app $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/macos/.
ssh $DEPLOY_USER@$DEPLOY_SERVER 'md5sum /var/www/html/dtac/agent/macos/dtac-agentd.app > /var/www/html/dtac/agent/macos/md5'
ssh $DEPLOY_USER@$DEPLOY_SERVER 'sha256sum /var/www/html/dtac/agent/macos/dtac-agentd.app > /var/www/html/dtac/agent/macos/sha256'
scp support/config/config.yaml $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/.
scp support/service/dtac-agentd.service $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/linux/.
scp support/service/dtac-agentd.service $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/macos/.
scp support/online-installers/windows/install.ps1 $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/windows/.
scp support/online-installers/linux/install.sh $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/linux/.
scp support/online-installers/darwin/install.sh $DEPLOY_USER@$DEPLOY_SERVER:/var/www/html/dtac/agent/macos/.