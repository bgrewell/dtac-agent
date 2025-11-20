# Webhost Modules

This directory contains webhost module implementations for DTAC.

## Overview

Webhost modules are DTAC-managed executables that serve embedded static web UIs and expose internal proxy routes. They run in two modes:

- **Standalone**: Local development with CLI/config and stdio logging
- **Managed**: Production deployment with gRPC control, token issuance, and event streaming

## Available Modules

### Dashboard (`dashboard/`)

Example webhost module demonstrating:
- Embedded static file serving
- Health check endpoint
- Proxy route configuration
- Both standalone and managed operation

## Quick Start

Build and run a module:

```bash
cd dashboard
go build -o dashboard
./dashboard
```

Access the dashboard at http://localhost:8080

## Documentation

- [Webhost Module Epic](../../docs/webhost-module-epic.md)
- [Developer Guide](../../docs/webhost-modules.md)

## Creating Your Own Module

See the developer guide at `docs/webhost-modules.md` for detailed instructions on creating custom webhost modules.
