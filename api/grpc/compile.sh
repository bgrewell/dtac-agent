#!/bin/bash

# Go
protoc --proto_path=. --go_out=paths=source_relative:go/. --go-grpc_out=paths=source_relative:go/. ./plugin.proto ./frontend.proto ./module.proto

# Python - There is a work around here due to Python having a broken protoc package/plugin/integration
# Check for the existence of the virtual environment directory
if [ ! -d ".venv" ]; then
    echo "Creating virtual environment..."
    python3 -m venv .venv
fi

# Activate the virtual environment
source .venv/bin/activate

# Install packages from requirements.txt if it exists
if [ -f "requirements.txt" ]; then
    echo "Installing packages from requirements.txt..."
    pip install -r requirements.txt
fi

# Run the protobuf compilation command
python3 -m grpc_tools.protoc --proto_path=. --python_out=./python --pyi_out=./python --grpc_python_out=./python plugin.proto module.proto

# Deactivate the virtual environment
deactivate

echo "Script execution completed."

#python3 -m grpc_tools.protoc --proto_path=. --python_out=./python --pyi_out=./python --grpc_python_out=./python plugin.proto
# TODO: DELETE ME python3 -m grpc_tools.protoc -I. --python_out=python/. --grpc_python_out=python/. plugin.proto
# TODO: DELETE ME protoc --proto_path=. plugin.proto --python_out=python/

# C# (Using Grpc.Tools NuGet package)
#protoc --csharp_out=dotnet/. --grpc_out=dotnet/. -I. --plugin=protoc-gen-grpc=grpc_csharp_plugin plugin.proto
