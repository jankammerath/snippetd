#!/bin/sh
# Build for arm64 linux 
echo "Building for arm64 linux..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o ./bin/snippetd ./main.go
echo "Building for arm64 linux done."

# Build the docker image
echo "Building the docker image..."
docker build -t snippetd .
echo "Building the docker image done."