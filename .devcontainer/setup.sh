#!/bin/bash
set -e

echo "🚀 Setting up Cribl SDK Sandbox..."

# Copy environment template
echo "📋 Copying .env template..."
cp .env.example .env

# Python setup
echo "🐍 Setting up Python dependencies..."
pip install --upgrade pip
pip install --upgrade --force-reinstall \
    git+https://github.com/criblio/cribl_control_plane_sdk_python.git \
    git+https://github.com/criblio/cribl_cloud_management_sdk_python.git
pip install -r requirements.txt

# Node.js setup
echo "📦 Setting up Node.js dependencies..."
npm install
npm update cribl-control-plane cribl-mgmt-plane
npm install -g typescript ts-node tsx

# Go setup - Control Plane
echo "🔧 Setting up Go Control Plane SDK..."
cd examples/control-plane/go
go get -u github.com/criblio/cribl-control-plane-sdk-go
go mod tidy

# Go setup - Management Plane
echo "🔧 Setting up Go Management Plane SDK..."
cd ../../mgmt-plane/go
go get -u github.com/criblio/cribl-cloud-management-sdk-go
go mod tidy

# Return to root
cd ../../..

echo "✅ Setup complete! All SDKs are ready to use."
