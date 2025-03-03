#!/bin/bash

# Generate Go code from Protocol Buffers
echo "Generating Go code from Protocol Buffers..."

# Generate proto for url.proto
protoc --go_out=. --go-grpc_out=. api/url.proto

# If you need to generate for all proto files, uncomment the line below
# protoc --go_out=. --go-grpc_out=. api/*.proto

echo "Code generation complete!"
