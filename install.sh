#!/bin/sh

set -e

# --- Check for root privileges ---
if [ "$EUID" -ne 0 ]; then
  echo "This script requires root privileges. Please run with sudo."
  exit 1
fi

# --- Check if binary exists ---
BINARY_PATH="./bin/snippetd"
if [ ! -f "$BINARY_PATH" ]; then
  echo "Error: Binary not found at $BINARY_PATH."

  echo "building the binary now..."
  ./build.sh
  if [ $? -ne 0 ]; then
    echo "Error: Failed to build the binary."
    exit 1
  fi

  echo "Binary built successfully. Continuing installation..."
fi

# --- Copy binary to /bin ---
echo "Installing binary to /bin/snippetd"
cp "$BINARY_PATH" /bin/snippetd
chmod +x /bin/snippetd

# --- Copy service file ---
echo "Installing systemd service file"
SERVICE_FILE="./snippetd.service"
if [ ! -f "$SERVICE_FILE" ]; then
  echo "Error: Service file not found at $SERVICE_FILE"
  exit 1
fi
cp "$SERVICE_FILE" /etc/systemd/system/snippetd.service

# --- Create snippetd user if it doesn't exist ---
USERNAME="snippetd"
ID_CMD=$(command -v id)
if ! $ID_CMD -u "$USERNAME" >/dev/null 2>&1; then
  echo "Creating user $USERNAME"
  useradd -r -m "$USERNAME"
fi

# --- Check if containerd group exists ---
GROUPNAME="containerd"
if ! getent group "$GROUPNAME" >/dev/null; then
  echo "Creating group $GROUPNAME"
  groupadd "$GROUPNAME"
fi

# --- Add snippetd user to containerd group ---
echo "Adding user $USERNAME to group $GROUPNAME"
usermod -a -G "$GROUPNAME" "$USERNAME"

# --- Change group of containerd socket ---
SOCKET_PATH="/run/containerd/containerd.sock"
if [ -S "$SOCKET_PATH" ]; then
  echo "Changing group of $SOCKET_PATH to $GROUPNAME"
  chgrp "$GROUPNAME" "$SOCKET_PATH"
  chmod g+rw "$SOCKET_PATH"
else
  echo "Warning: Socket file $SOCKET_PATH not found.  Containerd may not be running."
fi

# --- Register and launch service ---
echo "Registering and starting service"
systemctl daemon-reload
systemctl enable snippetd.service
systemctl start snippetd.service

echo "Installation complete!"