#!/usr/bin/env node
/**
 * Simple example: List Workspaces from Cribl Management Plane
 */

import { CriblMgmtPlane } from 'cribl-mgmt-plane';
import dotenv from 'dotenv';

// Load environment variables
dotenv.config();

async function listWorkspaces(): Promise<void> {
    console.log('Listing Cribl Workspaces');
    console.log('-'.repeat(40));

    // Get credentials from environment with placeholders
    const orgId = process.env.CRIBL_ORG_ID || 'your-org-id';
    const clientId = process.env.CRIBL_CLIENT_ID || 'your-client-id';
    const clientSecret = process.env.CRIBL_CLIENT_SECRET || 'your-client-secret';
    const domain = process.env.CRIBL_DOMAIN || 'cribl.cloud';

    // Check if credentials are properly set
    if (!orgId || !clientId || !clientSecret || 
        [orgId, clientId, clientSecret].some(val => val.startsWith('your-'))) {
        console.log('Missing credentials! Set these environment variables:');
        console.log('   CRIBL_ORG_ID');
        console.log('   CRIBL_CLIENT_ID');
        console.log('   CRIBL_CLIENT_SECRET');
        console.log('\nCopy .env.example to .env and fill in your values');
        return;
    }

    try {
        // Create client
        const client = new CriblMgmtPlane({
            serverURL: `https://gateway.${domain}`,
            security: {
                clientOauth: {
                    clientID: clientId,
                    clientSecret: clientSecret,
                    tokenURL: `https://login.${domain}/oauth/token`,
                    audience: `https://api.${domain}`
                }
            }
        });

        // List workspaces
        console.log('Fetching workspaces...');
        const response = await client.workspaces.list({ organizationId: orgId });
        
        // Handle the case where items might be undefined or empty
        const items = response.items || [];
        
        if (items.length > 0) {
            console.log(`\nFound ${items.length} workspace(s):`);
            console.log();
            
            for (const workspace of items) {
                const workspaceId = workspace.workspaceId || 'Unknown';
                console.log(`Workspace: ${workspaceId}`);
                console.log('-'.repeat(workspaceId.length + 12));
                
                // Print all available fields
                Object.entries(workspace).forEach(([key, value]) => {
                    console.log(`   ${key}: ${value}`);
                });
                console.log();
            }
        } else {
            console.log('No workspaces found');
        }

    } catch (error) {
        console.log(`Error: ${error}`);
    }
}

// Run the example
listWorkspaces().catch(console.error);
