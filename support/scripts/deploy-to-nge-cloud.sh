#!/usr/bin/env bash

make build
scp bin/system-agentd.exe intel@10.108.2.61:/var/www/html/systemapi/windows/.
scp bin/system-agentd intel@10.108.2.61:/var/www/html/systemapi/linux/.
scp support/config/config.yaml intel@10.108.2.61:/var/www/html/systemapi/.
scp support/service/system-agentd.service intel@10.108.2.61:/var/www/html/systemapi/linux/system-agentd.service