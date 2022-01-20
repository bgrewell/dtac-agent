#!/usr/bin/env bash

CONTAINER=protocbuilder
TAG=1.0

if [ ! -d go ]; then
  echo "[!] Creating go output directory"
  mkdir go;
fi

echo "[+] Building docker container"
docker image build -t $CONTAINER:$TAG .
docker container run --detach --name grpc $CONTAINER:$TAG
docker cp grpc:/go/src/github.com/BGrewell/system-api/plugin/go/plugin-api.pb.go ./go/.
#$(ls ./go/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}')  # strip the `omitempty` json attributes off of structures
echo "[+] Updating of go library complete"

echo "[+] Removing docker container"
docker rm grpc

echo "[+] Adding new files to source control"
git add go/plugin-api.pb.go
git commit -m "regenerated grpc libraries"

echo "[+] Done. Everything has been rebuilt and the repository has been updated"