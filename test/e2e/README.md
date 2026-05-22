# `/test/e2e`

End-to-end tests written as [DART](https://github.com/bgrewell/dart) workflows.
Each subdirectory contains one YAML suite plus any fixtures it needs.

DART manages container/VM lifecycle, so most suites run in a fresh LXD
container — see the workflow file's `nodes:` block for image choice and
network setup.

## Prerequisites

- DART CLI on `PATH`: `curl -sSL https://raw.githubusercontent.com/bgrewell/dart/main/install.sh | bash`
- LXD or Incus available locally, with the invoking user in the `lxd` group
- Any per-suite preconditions noted in the workflow header (e.g. an artifact
  produced by `mage osqueryBundle`)

## Suites

| Suite | Path | What it covers |
|---|---|---|
| osquery plugin | `osquery-plugin/osquery-plugin.yaml` | Spins up an Ubuntu 24.04 container, pushes the embedded `osquery.plugin`, validates `/health`, the curated read-only endpoints, generic `/query`, content-addressed extraction of `osqueryd`, and SIGTERM cleanup of the daemon child. |

## Running

From the repo root:

```sh
# Direct invocation
dart -c test/e2e/osquery-plugin/osquery-plugin.yaml -v

# Through mage (re-checks prerequisites, forwards $DART_FLAGS)
mage e2eOsquery
DART_FLAGS=-v mage e2eOsquery

# Setup / teardown only (useful when iterating)
dart -c test/e2e/osquery-plugin/osquery-plugin.yaml -setup
dart -c test/e2e/osquery-plugin/osquery-plugin.yaml -teardown
```

`dart` exits non-zero on the first failing test (with `-s`) or after running
the whole suite. The teardown step always runs, including when tests fail,
so the container is cleaned up automatically.
