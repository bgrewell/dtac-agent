#!/usr/bin/env bash

make build
scp bin/dtac-agentd.exe intel@10.108.2.61:/var/www/html/systemapi/windows/.
scp bin/dtac-agentd intel@10.108.2.61:/var/www/html/systemapi/linux/.
scp support/config/config.yaml intel@10.108.2.61:/var/www/html/systemapi/.
scp support/service/dtac-agentd.service intel@10.108.2.61:/var/www/html/systemapi/linux/dtac-agentd.service