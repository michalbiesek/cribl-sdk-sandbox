#!/bin/bash
set -e

echo "🚀 Setting up Cribl SDK Sandbox..."

# Copy environment template
echo "📋 Copying .env template..."
cp .env.example .env

# Python setup
echo "🐍 Setting up Python dependencies..."
pip install --upgrade pip
pip install -r requirements.txt

# Node.js setup
echo "📦 Setting up Node.js dependencies..."
npm install
npm install -g typescript ts-node tsx

# Go setup - Control Plane
echo "🔧 Setting up Go Control Plane SDK..."
cd examples/control-plane/go
go mod download
go mod tidy

# Go setup - Management Plane
echo "🔧 Setting up Go Management Plane SDK..."
cd ../../mgmt-plane/go
go mod download
go mod tidy

# Return to root
cd ../../..

echo "✅ Setup complete! All SDKs are ready to use."
