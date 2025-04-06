#!/bin/bash

# Check if the script is running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
  echo "This can only run on macOS (darwin/arm64)."
  exit 1
fi

# Check if LinuxKit is installed
if ! command -v linuxkit >/dev/null 2>&1; then
  echo "LinuxKit is not installed. Please install it."
  exit 1
fi

# Run the LinuxKit image with Apple Virtualization framework
echo "Running the LinuxKit image..."
linuxkit run virtualization linuxkit/linuxkit
echo "Running the LinuxKit image done."