#!/usr/bin/env bash

systems=( 10.108.1.21 10.108.2.61 10.108.1.11 )

for system in "${systems[@]}"
do
  echo "deploying to $system"
done