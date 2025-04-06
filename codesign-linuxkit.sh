#!/bin/sh
if [ "$EUID" -ne 0 ]; then
  echo "This script requires root privileges. Please run with sudo."
  exit 1
fi

LINUXKIT_EXEC=$(which linuxkit)
codesign --entitlements vz.entitlements -s - ${LINUXKIT_EXEC}