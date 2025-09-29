#!/usr/bin/env python3
"""
Simple example: List Worker Groups from On-Premise Cribl Control Plane
"""

import asyncio
import os
import logging
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
        
        # First, create an unauthenticated client to get a token
        client = CriblControlPlane(server_url=base_url)
        
        # Authenticate with username/password to get token
        print("🔐 Authenticating with username/password...")
        response = await client.auth.tokens.get_async(
            username=username, password=password
        )
        token = response.token
        print(f"✅ Authenticated with on-prem server, token: {token}")

        # Create authenticated SDK client with bearer token
        security = Security(bearer_auth=token)
        client = CriblControlPlane(server_url=base_url, security=security)
        client = CriblControlPlane(server_url=base_url, security=security)
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
        print(f"❌ Error: {error}")

if __name__ == "__main__":
    asyncio.run(list_worker_groups())
