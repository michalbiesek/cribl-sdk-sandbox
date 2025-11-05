#!/usr/bin/env python3
"""
Simple example: List Worker Groups from On-Premise Cribl Control Plane
"""

import asyncio
import os
import logging
import httpx
from dotenv import load_dotenv
from cribl_control_plane import CriblControlPlane
from cribl_control_plane.models import Security, ProductsCore

# Load environment variables
load_dotenv()

# Suppress verbose HTTP debug logging
logging.getLogger("httpcore").setLevel(logging.WARNING)
logging.getLogger("httpx").setLevel(logging.WARNING)

async def list_worker_groups():
    """List all worker groups in on-premise Cribl Control Plane."""
    print("🚀 Listing On-Premise Cribl Worker Groups")
    print("-" * 45)

    # Get credentials from environment with placeholders
    server_url = os.getenv("CRIBL_SERVER_URL") or "http://localhost:19000"
    username = os.getenv("CRIBL_USERNAME") or "admin"
    password = os.getenv("CRIBL_PASSWORD") or "admin"
    # Default to true for on-prem development environments with self-signed certs
    insecure_tls = os.getenv("CRIBL_INSECURE_TLS", "true").lower() != "false"

    # Show warning if using insecure TLS
    if server_url.startswith("https") and insecure_tls:
        print("⚠️  Accepting self-signed certificates (insecure mode)")

    # Check if server URL is properly set
    if server_url.startswith("your-"):
        print("❌ Invalid server URL! Set this environment variable:")
        print("   CRIBL_SERVER_URL")
        print("\n💡 Copy .env.example to .env and fill in your values")
        return

    try:
        # Create base URL for API
        base_url = f"{server_url.rstrip('/')}/api/v1"
        print(f"📡 Connecting to: {base_url}")
        
        # Configure custom HTTP client for HTTPS with self-signed certificates
        http_client = None
        if server_url.startswith("https") and insecure_tls:
            print("⚠️  Accepting self-signed certificates (insecure mode)")
            http_client = httpx.AsyncClient(verify=False)
        
        # First, create an unauthenticated client to get a token
        client = CriblControlPlane(
            server_url=base_url,
            async_client=http_client
        )
        
        # Authenticate with username/password to get token
        print("🔐 Authenticating with username/password...")
        response = await client.auth.tokens.get_async(username=username, password=password)
        token = response.token
        print(f"✅ Authenticated with on-prem server, token: {token}")

        # Create authenticated SDK client with bearer token
        client = CriblControlPlane(
            server_url=base_url,
            security=Security(bearer_auth=token),
            async_client=http_client
        )
        print("✅ Cribl SDK client created for on-prem server")

        # List worker groups
        print("📡 Fetching worker groups...")
        response = await client.groups.list_async(product=ProductsCore.STREAM)
        
        # Handle the case where items might be None or empty
        items = response.items or []
        
        if items:
            print(f"\n✅ Found {len(items)} worker group(s):")
            print()
            
            for group in items:
                print(f"📋 Worker Group: {group.id}")
                print("-" * (len(group.id) + 16))
                
                # Print all available fields
                for attr_name, value in vars(group).items():
                    print(f"   {attr_name}: {value}")
                print()
        else:
            print("📝 No worker groups found")

    except Exception as error:
        error_msg = str(error)
        print(f"❌ Error: {error_msg}")
        
        # Check if error is related to self-signed certificates
        if any(keyword in error_msg.lower() for keyword in [
            "certificate", "ssl", "self-signed", "cert_", "verify failed"
        ]):
            print("\n💡 Tip: If you're using a self-signed certificate, set:")
            print("   CRIBL_INSECURE_TLS=true")
            print("   (Only use this in development/testing environments!)")

if __name__ == "__main__":
    asyncio.run(list_worker_groups())
