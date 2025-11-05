# cribl-sdk-sandbox

Multi-language sandbox for experimenting with Cribl Control Plane and Cloud Management SDKs

## Supported SDKs

This sandbox provides examples and development environment support for the following Cribl SDKs:

### Control Plane SDKs

- **Python**: [cribl_control_plane_sdk_python](https://github.com/criblio/cribl_control_plane_sdk_python)
- **TypeScript**: [cribl-control-plane-sdk-typescript](https://github.com/criblio/cribl-control-plane-sdk-typescript)
- **Go**: [cribl-control-plane-sdk-go](https://github.com/criblio/cribl-control-plane-sdk-go)

### Management Plane SDKs

- **Python**: [cribl_cloud_management_sdk_python](https://github.com/criblio/cribl_cloud_management_sdk_python)
- **TypeScript**: [cribl-cloud-management-sdk-typescript](https://github.com/criblio/cribl-cloud-management-sdk-typescript)
- **Go**: [cribl-cloud-management-sdk-go](https://github.com/criblio/cribl-cloud-management-sdk-go)

## Quick Start

1. **Open in devcontainer** - All dependencies are automatically set up
2. **Configure credentials** - Edit `.env` with your Cribl credentials
3. **Run examples** - Use Python, TypeScript, or Go examples

### Cloud

Connect to Cribl.Cloud using OAuth2 authentication:

- `CRIBL_ORG_ID` - Your organization ID
- `CRIBL_CLIENT_ID` - OAuth2 client ID
- `CRIBL_CLIENT_SECRET` - OAuth2 client secret
- `CRIBL_DOMAIN` - Domain for Cribl Cloud (defaults to `cribl.cloud`)
- `CRIBL_WORKSPACE_NAME` - Workspace name (defaults to `main`)

For detailed authentication setup, see: [Cribl API Documentation](https://docs.cribl.io/api/#criblcloud)

### On-Premise

Connect to an on-premise Cribl Control Plane:

- `CRIBL_SERVER_URL` - Server URL (e.g., `http://localhost:19000`)
- `CRIBL_USERNAME` - Username (defaults to `admin`)
- `CRIBL_PASSWORD` - Password (defaults to `admin`)

#### Local Testing

Start a local Cribl Stream leader for testing:

```bash
./start-cribl.sh
```

This provides a local environment at http://localhost:19000 with default credentials `admin/admin`.

To stop:

```bash
./stop-cribl.sh
```

## Examples

Run examples in Python, TypeScript, or Go:

```bash
# Control Plane Examples
python examples/control-plane/python/example_cloud.py
python examples/control-plane/python/example_onprem.py

# Management Plane Examples
python examples/mgmt-plane/python/example.py

# Or use npm scripts
npm run control-plane:cloud
npm run control-plane:onprem
npm run mgmt-plane:example

# Or run Go examples directly
cd examples/control-plane/go && go run example_cloud.go
cd examples/mgmt-plane/go && go run example.go
```

## VS Code Integration

VS Code launch configurations and tasks are provided for running and debugging examples.

## Ports

- **19000** - Cribl Stream UI
- **4200** - Cribl leader communication port
