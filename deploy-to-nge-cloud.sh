#!/usr/bin/env bash

make build
scp bin/system-apid.exe intel@10.108.2.61:/var/www/html/systemapi/windows/.
scp bin/system-apid intel@10.108.2.61:/var/www/html/systemapi/linux/.
scp support/config/config.yaml intel@10.108.2.61:/var/www/html/systemapi/.
scp support/service/system-apid.service intel@10.108.2.61:/var/www/html/systemapi/linux/system-apid.service