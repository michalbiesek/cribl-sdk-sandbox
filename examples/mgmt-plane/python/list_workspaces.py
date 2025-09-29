#!/usr/bin/env python3
"""
Simple example: List Workspaces from Cribl Management Plane
"""

import asyncio
import os
import logging
from dotenv import load_dotenv
from cribl_mgmt_plane import CriblMgmtPlane
from cribl_mgmt_plane.models import Security, SchemeClientOauth

# Load environment variables
load_dotenv()

# Suppress verbose HTTP debug logging
logging.getLogger("httpcore").setLevel(logging.WARNING)
logging.getLogger("httpx").setLevel(logging.WARNING)

async def list_workspaces():
    """List all workspaces in Cribl Management Plane."""
    print("🚀 Listing Cribl Workspaces")
    print("-" * 40)

    # Get credentials from environment with placeholders
    org_id = os.getenv("CRIBL_ORG_ID") or "your-org-id"
    client_id = os.getenv("CRIBL_CLIENT_ID") or "your-client-id"
    client_secret = os.getenv("CRIBL_CLIENT_SECRET") or "your-client-secret"

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
            token_url="https://login.cribl.cloud/oauth/token",
            audience="https://api.cribl.cloud"
        )

        # Create client  
        base_url = f"https://gateway.cribl.cloud"
        security = Security(client_oauth=client_oauth)
        client = CriblMgmtPlane(server_url=base_url, security=security)

        # List workspaces
        print("📡 Fetching workspaces...")
        response = await client.workspaces.list_async(organization_id=org_id)
        
        # Handle the case where items might be None or empty
        items = response.items or []
        
        if items:
            print(f"\n✅ Found {len(items)} workspace(s):")
            print()
            
            for workspace in items:
                # Use the correct workspace identifier from the schema
                workspace_id = getattr(workspace, 'workspaceId', 'Unknown')
                print(f"🏢 Workspace: {workspace_id}")
                print("-" * (len(str(workspace_id)) + 12))
                
                # Print all available fields
                for attr_name, value in vars(workspace).items():
                    print(f"   {attr_name}: {value}")
                print()
        else:
            print("📝 No workspaces found")

    except Exception as error:
        print(f"❌ Error: {error}")

if __name__ == "__main__":
    asyncio.run(list_workspaces())
