#!/bin/bash

# Exit on error
set -e

# Configuration
BUILD_DIR="bin"
CONFIG_DIR="config"
LOG_DIR="logs"

# Create necessary directories
mkdir -p $BUILD_DIR
mkdir -p $LOG_DIR

# Build server
echo "Building server..."
go build -o $BUILD_DIR/server ./server

# Build client
echo "Building client..."
go build -o $BUILD_DIR/client ./client

# Copy configuration files
echo "Copying configuration files..."
cp $CONFIG_DIR/*.json $BUILD_DIR/

# Set permissions
chmod +x $BUILD_DIR/server
chmod +x $BUILD_DIR/client

echo "Build completed successfully!"
echo "Server binary: $BUILD_DIR/server"
echo "Client binary: $BUILD_DIR/client" 