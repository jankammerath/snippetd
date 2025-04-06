#!/bin/bash

# Check if the script is running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
  echo "This script is intended to be run on macOS."
  exit 1
fi

# Build for arm64 linux 
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/snippetd ./main.go