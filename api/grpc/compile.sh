#!/bin/bash

# Go
protoc --proto_path=. plugin.proto --go_out=paths=source_relative:go/. --go-grpc_out=paths=source_relative:go/.

# Python
python3 -m grpc_tools.protoc -I. --python_out=python/. --grpc_python_out=python/. plugin.proto

# C# (Using Grpc.Tools NuGet package)
#protoc --csharp_out=dotnet/. --grpc_out=dotnet/. -I. --plugin=protoc-gen-grpc=grpc_csharp_plugin plugin.proto
