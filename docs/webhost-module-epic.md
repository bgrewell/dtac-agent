# Design and implement webhost module infrastructure (separate from plugins)

Overview

Implement support for a new class of "webhost modules" in DTAC, distinct from current backend plugins. These modules will serve embedded static web assets and expose HTTP routes (including user-defined proxy routes that forward to internal services the browser must not access directly). Modules operate in two modes: DTAC-managed (gRPC control plane) and standalone (local dev/test). This epic captures design, configuration, subsystem, token/Casbin integration, developer SDK, example module, and documentation.

Goals

- Provide a gRPC control API (WebhostService) to initialize modules, issue tokens, and stream logs/events.
- Ensure DTAC config fully controls listen address/port, proxy definitions, and permission scopes for each module.
- Provide a Go SDK to minimize module boilerplate (mode detection, logging abstraction, token helper, proxy helper).
- Allow modules to run standalone for local development with equivalent behavior but local config and stdio logging.

Core tasks (high-level)

1. Design and document webhost module architecture and gRPC API
2. Extend DTAC configuration model with a `webhosts:` section
3. Implement the `internal/webhost` subsystem to manage module lifecycle
4. Implement token issuance (JWTs) for modules with Casbin enforcement
5. Provide `pkg/webhost` SDK (base types, host, proxy helpers)
6. Add an example module under `cmd/webhosts/dashboard`
7. Integrate logs/events into DTAC logging and add developer documentation

Acceptance criteria

- Draft design doc and gRPC proto in the repo/docs
- DTAC can spawn and initialize a webhost module (managed mode)
- Example module can run standalone and under DTAC-managed mode
- Token requests are enforced by Casbin and return signed JWTs when allowed
- Developer guide shows how to author and test modules

This file is the epic description and contains links to the starter code included in this PR.
