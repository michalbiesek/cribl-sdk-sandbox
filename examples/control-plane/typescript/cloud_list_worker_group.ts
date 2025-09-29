#!/usr/bin/env node
/**
 * Simple example: List Worker Groups from Cribl Cloud Control Plane
 */

import dotenv from 'dotenv';
import { CriblControlPlane } from 'cribl-control-plane';

// Load environment variables
dotenv.config();

async function listWorkerGroups(): Promise<void> {
    console.log('Listing Cribl Worker Groups');
    console.log('-'.repeat(40));

    // Get credentials from environment with placeholders
    const orgId = process.env.CRIBL_ORG_ID || 'your-org-id';
    const clientId = process.env.CRIBL_CLIENT_ID || 'your-client-id';
    const clientSecret = process.env.CRIBL_CLIENT_SECRET || 'your-client-secret';
    const workspace = process.env.CRIBL_WORKSPACE_NAME || 'main';

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
        const client = new CriblControlPlane({
            serverURL: `https://${workspace}-${orgId}.cribl.cloud/api/v1`,
            security: {
                clientOauth: {
                    clientID: clientId,
                    clientSecret: clientSecret,
                    tokenURL: 'https://login.cribl.cloud/oauth/token',
                    audience: 'https://api.cribl.cloud'
                }
            }
        });

        // List worker groups  
        console.log('Fetching worker groups...');
        const response = await client.groups.list({ product: 'stream' });
        
        // Handle the case where items might be undefined or empty
        const items = response.items || [];
        
        if (items.length > 0) {
            console.log(`\nFound ${items.length} worker group(s):`);
            console.log();
            
            for (const group of items) {
                const groupId = group.id || 'Unknown';
                console.log(`Worker Group: ${groupId}`);
                console.log('-'.repeat(groupId.length + 16));
                
                // Print all available fields
                Object.entries(group).forEach(([key, value]) => {
                    console.log(`   ${key}: ${value}`);
                });
                console.log();
            }
        } else {
            console.log('No worker groups found');
        }

    } catch (error) {
        console.log(`Error: ${error}`);
    }
}

// Run the example
listWorkerGroups().catch(console.error);
