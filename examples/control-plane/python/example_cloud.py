#!/usr/bin/env python3
"""
Simple example: List Worker Groups from Cribl Cloud
"""

import asyncio
import os
import logging
from dotenv import load_dotenv
from cribl_control_plane import CriblControlPlane
from cribl_control_plane.models import Security, SchemeClientOauth, ProductsCore

# Load environment variables
load_dotenv()

# Suppress verbose HTTP debug logging
logging.getLogger("httpcore").setLevel(logging.WARNING)
logging.getLogger("httpx").setLevel(logging.WARNING)

async def list_worker_groups():
    """List all worker groups in Cribl Cloud."""
    print("🚀 Listing Cribl Worker Groups")
    print("-" * 40)

    # Get credentials from environment with placeholders
    org_id = os.getenv("CRIBL_ORG_ID") or "your-org-id"
    client_id = os.getenv("CRIBL_CLIENT_ID") or "your-client-id"
    client_secret = os.getenv("CRIBL_CLIENT_SECRET") or "your-client-secret"
    workspace = os.getenv("CRIBL_WORKSPACE_NAME") or "main"
    domain = os.getenv("CRIBL_DOMAIN") or "cribl.cloud"

    # Check if credentials are properly set (not empty or placeholders)
    if not all([org_id, client_id, client_secret]) or any(val.startswith("your-") for val in [org_id, client_id, client_secret]):
        print("❌ Missing credentials! Set these environment variables:")
        print("   CRIBL_ORG_ID")
        print("   CRIBL_CLIENT_ID") 
        print("   CRIBL_CLIENT_SECRET")
        print("\n💡 Copy .env.example to .env and fill in your values")
        return

    try:
        # Setup authentication
        client_oauth = SchemeClientOauth(
            client_id=client_id,
            client_secret=client_secret,
            token_url=f"https://login.{domain}/oauth/token",
            audience=f"https://api.{domain}"
        )

        # Create client
        base_url = f"https://{workspace}-{org_id}.{domain}/api/v1"
        security = Security(client_oauth=client_oauth)
        client = CriblControlPlane(server_url=base_url, security=security)

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

