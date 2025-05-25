#!/bin/bash

# Exit on error
set -e

# Configuration
REMOTE_USER="your-username"
REMOTE_HOST="your-server"
REMOTE_DIR="/opt/tcr"
LOCAL_BUILD_DIR="bin"

# Build the project first
./scripts/build.sh

# Create remote directory if it doesn't exist
ssh $REMOTE_USER@$REMOTE_HOST "mkdir -p $REMOTE_DIR"

# Copy files to remote server
echo "Copying files to remote server..."
scp -r $LOCAL_BUILD_DIR/* $REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/

# Set up systemd service
echo "Setting up systemd service..."
cat > tcr.service << EOL
[Unit]
Description=Text-Based Clash Royale Server
After=network.target

[Service]
Type=simple
User=$REMOTE_USER
WorkingDirectory=$REMOTE_DIR
ExecStart=$REMOTE_DIR/server -config $REMOTE_DIR/prod.json
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOL

# Copy service file and enable it
scp tcr.service $REMOTE_USER@$REMOTE_HOST:/tmp/
ssh $REMOTE_USER@$REMOTE_HOST "sudo mv /tmp/tcr.service /etc/systemd/system/ && \
    sudo systemctl daemon-reload && \
    sudo systemctl enable tcr && \
    sudo systemctl restart tcr"

# Clean up
rm tcr.service

echo "Deployment completed successfully!"
echo "Server is running at $REMOTE_HOST" 