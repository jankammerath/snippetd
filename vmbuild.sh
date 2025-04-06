#!/bin/bash

# Check if the script is running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
  echo "This script is intended to be run on macOS."
  exit 1
fi

# Check if LinuxKit is installed
if ! command -v linuxkit >/dev/null 2>&1; then
  echo "linuxkit is not installed. Please install it."
  exit 1
fi

# Build for arm64 linux 
echo "Building for arm64 linux..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/snippetd ./main.go
echo "Building for arm64 linux done."

# Build the LinuxKit image into the vm folder
echo "Building the LinuxKit image..."
linuxkit build --dir ./vm linuxkit.yml
echo "Building the LinuxKit image done."