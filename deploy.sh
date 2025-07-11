#!/bin/bash
set -euo pipefail

# Binary names and target hosts
CLIENT_BIN=netlsrc
SERVER_BIN=netlsrd
SERVER_HOST=192.168.2.127
CLIENT_HOST=192.168.2.128
USER=mo

# Build the Go project binaries
echo "Building $SERVER_BIN for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o "$SERVER_BIN" ./cmd/netlsrd

echo "Building $CLIENT_BIN for linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o "$CLIENT_BIN" ./cmd/netlsrc

# Deploy server binary and config to server host
echo "Deploying server to $SERVER_HOST..."
rsync -avz --progress "$SERVER_BIN" "${USER}@${SERVER_HOST}:~/"
rsync -avz --progress "configs/server.yaml" "${USER}@${SERVER_HOST}:~/config.yaml"

# Deploy client binary and config to client host
echo "Deploying client to $CLIENT_HOST..."
rsync -avz --progress "$CLIENT_BIN" "${USER}@${CLIENT_HOST}:~/"
rsync -avz --progress "configs/client.yaml" "${USER}@${CLIENT_HOST}:~/config.yaml"

# Start server first
echo "Starting server on $SERVER_HOST..."
ssh "${USER}@${SERVER_HOST}" "sudo pkill -x $SERVER_BIN || true"
ssh "${USER}@${SERVER_HOST}" "nohup /home/${USER}/$SERVER_BIN -config config.yaml > netlsr.log 2>&1 &"

# Wait a moment for server to start
sleep 3

# Start client
echo "Starting client on $CLIENT_HOST..."
ssh "${USER}@${CLIENT_HOST}" "sudo pkill -x $CLIENT_BIN || true"
ssh "${USER}@${CLIENT_HOST}" "nohup /home/${USER}/$CLIENT_BIN -config config.yaml > netlsr.log 2>&1 &"

echo "Deployment complete."
echo "Server running on $SERVER_HOST"
echo "Client running on $CLIENT_HOST" 


