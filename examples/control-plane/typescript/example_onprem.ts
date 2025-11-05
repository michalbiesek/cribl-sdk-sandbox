#!/usr/bin/env node
/**
 * Simple example: List Worker Groups from On-Premise Cribl Control Plane
 */

import { CriblControlPlane } from "cribl-control-plane";
import dotenv from "dotenv";

// Load environment variables
dotenv.config();

async function listWorkerGroups(): Promise<void> {
  console.log("Listing On-Premise Cribl Worker Groups");
  console.log("-".repeat(45));

  // Get credentials from environment with placeholders
  const serverUrl = process.env.CRIBL_SERVER_URL || "http://localhost:19000";
  const username = process.env.CRIBL_USERNAME || "admin";
  const password = process.env.CRIBL_PASSWORD || "admin";
  // Default to true for on-prem development environments with self-signed certs
  const insecureTls = process.env.CRIBL_INSECURE_TLS !== "false";

  // Check if server URL is properly set
  if (serverUrl.startsWith("your-")) {
    console.log("Invalid server URL! Set this environment variable:");
    console.log("   CRIBL_SERVER_URL");
    console.log(
      "\nCopy examples/.env.example to examples/.env and fill in your values"
    );
    return;
  }

  try {
    // Create base URL for API
    const baseUrl = `${serverUrl.replace(/\/$/, "")}/api/v1`;
    console.log(`Connecting to: ${baseUrl}`);

    // Configure TLS for HTTPS with self-signed certificates
    if (serverUrl.startsWith("https") && insecureTls) {
      console.log("⚠️  Accepting self-signed certificates (insecure mode)");
      process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
    }

    // First, create an unauthenticated client to get a token
    let client = new CriblControlPlane({ serverURL: baseUrl });

    // Authenticate with username/password to get token
    console.log("Authenticating with username/password...");
    const authResponse = await client.auth.tokens.get({ username, password });
    const token = authResponse.token;
    console.log(`Authenticated with on-prem server, token: ${token}`);

    // Create authenticated SDK client with bearer token
    client = new CriblControlPlane({
      serverURL: baseUrl,
      security: { bearerAuth: token },
    });
    console.log("Cribl SDK client created for on-prem server");

    // List worker groups
    console.log("Fetching worker groups...");
    const response = await client.groups.list({ product: "stream" });

    // Handle the case where items might be undefined or empty
    const items = response.items || [];

    if (items.length > 0) {
      console.log(`\nFound ${items.length} worker group(s):`);
      console.log();

      for (const group of items) {
        const groupId = group.id || "Unknown";
        console.log(`Worker Group: ${groupId}`);
        console.log("-".repeat(groupId.length + 16));

        // Print all available fields
        Object.entries(group).forEach(([key, value]) => {
          console.log(`   ${key}: ${value}`);
        });
        console.log();
      }
    } else {
      console.log("No worker groups found");
    }
  } catch (error) {
    const errorMsg = String(error);
    console.log(`Error: ${errorMsg}`);

    // Check if error is related to self-signed certificates
    if (
      errorMsg.includes("self signed certificate") ||
      errorMsg.includes("DEPTH_ZERO_SELF_SIGNED_CERT") ||
      errorMsg.includes("certificate") ||
      errorMsg.includes("CERT_")
    ) {
      console.log("\n💡 Tip: If you're using a self-signed certificate, set:");
      console.log("   CRIBL_INSECURE_TLS=true");
      console.log("   (Only use this in development/testing environments!)");
    }
  }
}

// Run the example
listWorkerGroups().catch(console.error);
