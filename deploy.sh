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
  if [ "$HOST" == "192.168.2.127" ]; then
    ssh "${USER}@${HOST}" "sudo pkill -x netlsr || true"
    ssh "${USER}@${HOST}" "nohup sudo netlsr -mode server -ifname tun77 -local-ip 10.177.0.1/24 -debug > netlsr.log 2>&1 &"
  else
    sleep 2
    ssh "${USER}@${HOST}" "sudo pkill -x netlsr || true"
    ssh "${USER}@${HOST}" "nohup sudo netlsr -mode client -remote 192.168.2.127 -ifname tun77 -local-ip 10.177.0.2/24 > netlsr.log 2>&1 &"
  fi
done

echo "Deployment complete." 