# Webhost Modules — Developer Guide

Overview

Webhost modules are a new class of DTAC-managed processes that host web UIs and optionally proxy API requests to internal services. They run in two modes:
- DTAC-managed: full control via gRPC (Init, token issuance, event/log stream).
- Standalone: local dev mode; module reads local config / flags and logs to stdout/stderr.

Quickstart (example)

1. Build the example module:
   go build ./cmd/webhosts/dashboard

2. Run standalone:
   ./dashboard --standalone --listen-addr 127.0.0.1 --port 3000

3. Configure DTAC to manage the module:
   Add a `webhosts` entry in dtac config with `executable` and `listen` info. DTAC will spawn the module with `DTAC_WEBHOSTS=true`, parse CONNECT{{...}}, connect via gRPC and call `Init`.

Key concepts

- DTAC owns listen address and proxy configuration in managed mode.
- Proxies are registered into the module via `Init` and the module should setup reverse proxies that inject credentials (JWTs) obtained via the `RequestToken` RPC.
- Logging in managed mode flows through the EventStream; standalone logs go to stderr/stdout.

Security

- Proxies keep internal service URLs and credentials off the browser.
- DTAC controls which scopes a module is allowed to request; tokens are signed by DTAC.
