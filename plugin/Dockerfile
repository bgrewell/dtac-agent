# Dockerfile.protogen
FROM golang:latest

LABEL maintainer="Benjamin Grewell <benjamin.grewell@intel.com>"

ENV PROTOC_VERSION 3.6.1
ENV PROTOC_GEN_GO_VERSION v1.2.0

WORKDIR /go/src/github.com/BGrewell/system-api/plugin/

RUN apt update
RUN apt install -y protobuf-compiler
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN export PATH=$PATH:$GOPATH/bin
RUN echo $PATH

COPY api api

RUN mkdir go
RUN protoc -I api/ --go_out=plugins=grpc:go --go_opt=paths=source_relative api/plugin-api.proto
