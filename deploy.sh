#!/usr/bin/env bash
set -euo pipefail

# Binary name and target hosts
BIN=netlsr
HOSTS=("192.168.2.127" "192.168.2.128")
USER=mo

# Build the Go project
echo "Building $BIN for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o "$BIN"

# Deploy to each host
for HOST in "${HOSTS[@]}"; do
  echo "Deploying to $HOST..."
  rsync -avz --progress "$BIN" "${USER}@${HOST}:~/"
done

echo "Deployment complete." 