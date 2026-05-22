# `osquery` plugin

Wraps a private `osqueryd` subprocess and exposes its read-only query surface
through the dtac plugin framework. The plugin owns the daemon's lifecycle:
spawn at startup, supervise, and reap on `SIGTERM` / `SIGINT`.

## Build modes

The plugin works in two distribution shapes; pick one at build time:

| Mode | Build command | What you ship | When to use |
|---|---|---|---|
| **Embedded (recommended)** | `mage osqueryBundle` | One `osquery-<os>-<arch>.plugin` binary (~300 MiB; embeds osqueryd 5.23.0). At first run, extracts osqueryd to a content-addressed cache directory. | Default. Zero-install, one file to deploy. |
| **External osqueryd** | `mage plugins` | A slim `osquery.plugin` (~30 MiB) + you supply `osqueryd` separately via `config.binary_path` or `osquery-bin/osqueryd` next to the plugin. | Air-gapped environments where you have your own vetted osqueryd, or where binary size matters. |

The build tag that gates the embed is `embed_osquery`; the `osqueryBundle`
mage target sets it for you.

Embedded-supported platforms today: `linux/amd64`, `linux/arm64`, `darwin/arm64`.
`darwin/amd64` has no native upstream tarball — it falls through to the
external-osqueryd path.

## Run modes

Mode is env-driven, so the same binary runs both ways without a rebuild:

```sh
# Embedded gRPC mode — the dtac-agent spawns and talks to it via gRPC.
# Default; no env vars needed.
/opt/dtac/plugins/osquery.plugin

# Standalone REST mode — the plugin is its own HTTPS service.
DTAC_STANDALONE=true DTAC_STANDALONE_PORT=8080 ./osquery.plugin
# Endpoints land at: http://host:8080/osquery/<endpoint>
# Health check at:   http://host:8080/health
```

Env vars (consumed by `pkg/plugins.NewStandaloneConfig`):
`DTAC_STANDALONE`, `DTAC_STANDALONE_PROTOCOL`, `DTAC_STANDALONE_PORT`,
`DTAC_STANDALONE_HOST`, `DTAC_STANDALONE_TLS_CERT`, `DTAC_STANDALONE_TLS_KEY`.

## Configuration

JSON config keys (under `plugins.entries.osquery.config` in the agent's
`config.yaml`):

| Key | Default | Meaning |
|---|---|---|
| `binary_path` | unset | Absolute path to an existing `osqueryd` — bypasses the embed/sibling lookup. Use this for system-installed osquery. |
| `extra_flags` | `[]` | Extra arguments appended to the `osqueryd` command line. |
| `socket_path` | `$TMPDIR/dtac-osquery-<pid>.sock` | Override the unix socket used by the Thrift extensions channel. |
| `query_timeout` | `10s` | Per-query hard timeout. |
| `startup_timeout` | `5s` | How long to wait for `osqueryd` to publish its extensions socket. |
| `max_rows` | `10000` | Safety cap on rows returned by `/query`. Truncated responses set `truncated: true`. |

## Endpoint surface

All endpoints under root `osquery/`. With dtac-agent: `/plugins/osquery/<path>`.
In standalone REST mode: `/osquery/<path>`. The response envelope is the same
across all read endpoints:

```json
{ "rows": [{...}, ...], "row_count": N, "duration_ms": N, "truncated": false }
```

### Generic

| Path | Method | Auth | Notes |
|---|---|---|---|
| `query` | POST | admin | Body `{"sql":"…"}`. Returns the standard envelope. SQL is opaque; osquery is read-only by design but you should still gate this at admin. |

### Identity / inventory (User tier — safe to expose broadly)

| Path | Method | Auth | Source table |
|---|---|---|---|
| `os_version` | GET | user | `os_version` |
| `system_info` | GET | user | `system_info` |
| `kernel_info` | GET | user | `kernel_info` |
| `uptime` | GET | user | `uptime` |
| `osquery_info` | GET | user | `osquery_info` |

### Operator-tier (network / topology / inventory)

| Path | Method | Auth | Source table(s) |
|---|---|---|---|
| `processes` | GET | operator | `processes` |
| `users` | GET | operator | `users` |
| `groups` | GET | operator | `groups` |
| `listening_ports` | GET | operator | `listening_ports` |
| `interface_addresses` | GET | operator | `interface_addresses` |
| `routes` | GET | operator | `routes` |
| `arp_cache` | GET | operator | `arp_cache` |
| `dns_resolvers` | GET | operator | `dns_resolvers` |
| `etc_hosts` | GET | operator | `etc_hosts` |
| `kernel_modules` | GET | operator | `kernel_modules` (Linux only) |
| `block_devices` | GET | operator | `block_devices` (Linux only) |
| `mounts` | GET | operator | `mounts` |
| `process_open_sockets` | GET | operator | `process_open_sockets` |
| `active_sessions` | GET | operator | `processes` ⋈ `process_open_sockets` ⋈ `users` (sshd children) |
| `logged_in_users` | GET | operator | `logged_in_users` (utmp) |
| `services` | GET | operator | `systemd_units` (Linux) / `launchd` (darwin) |
| `installed_packages` | GET | operator | `deb_packages` + `rpm_packages` (Linux) / `apps` (darwin) |

### Admin-tier (audit / persistence / secrets-adjacent)

| Path | Method | Auth | Source table |
|---|---|---|---|
| `crontab` | GET | admin | `crontab` |
| `last_logins` | GET | admin | `last` (wtmp) |
| `sudoers` | GET | admin | `sudoers` |

### `active_sessions` — why a separate endpoint from `logged_in_users`

`logged_in_users` and the traditional `who` / `w` / `users` utilities all
read `utmp`, which only records sessions that **allocated a TTY**. They miss
every TTY-less SSH activity:

- `ssh user@host 'cmd'` — runs the command and exits, no TTY, no utmp.
- `ssh -N -L 8080:internal:80 user@host` — port-forward only, never any shell.
- `ProxyJump` / `ProxyCommand` transits — the jump host's sshd holds the
  forwarded connection but never opens a session.
- `scp` / `sftp` — short-lived, no utmp.
- SSH connection multiplexing — only the master shows up; slaves don't.

`active_sessions` enumerates every authenticated SSH connection by walking
`sshd` child processes and joining them against `process_open_sockets` filtered
to established remote sockets. That invariant catches all of the cases above.

## Privilege requirements

The plugin runs as whatever user the dtac-agent (or the standalone process)
runs as. osquery's coverage degrades gracefully on partial privilege: tables
return what the current uid can see. Effective coverage by tier:

### Works fully as a non-root user

`os_version`, `system_info`, `kernel_info`, `uptime`, `osquery_info`,
`interface_addresses`, `routes`, `arp_cache`, `dns_resolvers`, `etc_hosts`,
`kernel_modules`, `mounts`, `groups`, `users`, `installed_packages`,
`logged_in_users`, `services`.

These all read world-readable kernel interfaces (`/proc`, netlink, sysfs)
or world-readable config files.

### Works as non-root but with degraded visibility

| Endpoint | Without root you see… |
|---|---|
| `processes` | All PIDs and names, but truncated `cmdline` / `exe` for processes owned by other users. |
| `listening_ports` | Only sockets owned by the current uid. |
| `process_open_sockets` | Only sockets owned by the current uid. |
| `active_sessions` | Only your own sshd children — typically nothing useful. |
| `block_devices` | Listed, but `uuid`/`label` may be empty without `/dev` read access. |
| `last_logins` | Works if the user is in the `utmp` group (varies by distro). |

### Requires root

| Endpoint | Why |
|---|---|
| `sudoers` | `/etc/sudoers` is `0440 root:root` (and `/etc/sudoers.d/*` similarly). |
| `crontab` | Non-root only sees the invoking user's own crontab; the system view needs root. |

### Capabilities-only alternative (Linux)

If you don't want to run the plugin as `root`, the minimum capability set
for full coverage is:

- `CAP_DAC_READ_SEARCH` — read `/etc/sudoers`, other users' crontabs, all of `/proc`.
- `CAP_SYS_PTRACE` — full `cmdline`/`exe` for other users' processes.
- `CAP_NET_ADMIN` — full visibility into `process_open_sockets` / `listening_ports`.
- `CAP_SYS_NICE` (optional) — read process priority / cpu affinity.

Granting via systemd:

```ini
[Service]
AmbientCapabilities=CAP_DAC_READ_SEARCH CAP_SYS_PTRACE CAP_NET_ADMIN
CapabilityBoundingSet=CAP_DAC_READ_SEARCH CAP_SYS_PTRACE CAP_NET_ADMIN
NoNewPrivileges=true
```

Or via `setcap` on the plugin binary:

```sh
sudo setcap 'cap_dac_read_search,cap_sys_ptrace,cap_net_admin+eip' \
    /opt/dtac/plugins/osquery.plugin
```

Note: the *embedded* osqueryd inherits the plugin's caps because the plugin
spawns it as a child. The caps on the plugin binary need to be ambient
(`+eip`), not just effective (`+ep`).

### Containers / namespaces

When the plugin runs inside an unprivileged container (LXC, Docker without
`--privileged`, Kubernetes with default `securityContext`), "root inside the
container" is mapped to an unprivileged uid on the host. The privilege model
above describes what `uid 0` inside the namespace sees — which is normally
all the in-namespace data you want. Endpoints that escape namespaces
(e.g. host-level `kernel_modules`) return what the container's view shows,
not the host's.

### SELinux / AppArmor

osquery's default SELinux/AppArmor profiles on RHEL/Ubuntu are written for
the upstream `/opt/osquery/bin/osqueryd` path — the *extracted* binary at
`/var/cache/dtac/osquery/<sha>/osqueryd` won't match those profiles. If you
run an enforcing system policy, either:

1. Use `config.binary_path` to point at the system-packaged osqueryd
   (and skip the embed), so the upstream profile applies, or
2. Author a custom profile for the cache path; or
3. Disable the embed and ship a system osquery package via your normal
   package mgmt.

## Runtime cache

When using the embedded mode, osqueryd is extracted on first run to:

1. `/var/cache/dtac/osquery/<sha256>/osqueryd` — preferred; survives reboots.
   Requires write access to `/var/cache/dtac` (root, or a directory
   pre-created with the right ownership).
2. `$TMPDIR/dtac-osquery/<sha256>/osqueryd` — fallback; lost on reboot.

The directory name is the sha256 of the embedded bytes, so plugin upgrades
naturally create new cache entries instead of overwriting in place.

To force re-extraction (e.g. for testing), remove the cache directory.

## Troubleshooting

- **`osqueryd extensions socket did not appear within 5s`** — osqueryd
  started but didn't open its socket. Bump `startup_timeout` or check
  `/tmp/plugin.log` (standalone mode) or the agent log (embedded mode) for
  osqueryd stderr.
- **`Permission denied` on `/var/cache/dtac/osquery`** — the plugin will
  fall back to `$TMPDIR` automatically; if both fail you'll see a clear
  extraction error. Pre-create the cache dir with appropriate ownership to
  pin it on `/var/cache`.
- **`/active_sessions` returns nothing on a host with active SSH sessions** —
  almost always a privilege issue. The plugin can only see its own user's
  sshd children. Run as root or grant `CAP_NET_ADMIN` + `CAP_SYS_PTRACE`.
- **`/sudoers` returns empty** — the plugin lacks read access to `/etc/sudoers`.
  Run as root or grant `CAP_DAC_READ_SEARCH`.

## Source layout

```
cmd/plugins/osquery/
├── main.go                       # entry point; signal-driven shutdown
├── build.yaml                    # picked up by `mage plugins`
└── osqueryplugin/
    ├── osquery_plugin.go         # OsqueryPlugin (Register, runQueryWrapped, Shutdown)
    ├── config.go                 # JSON config + defaults
    ├── daemon.go                 # osqueryd subprocess supervision
    ├── extractor.go              # content-addressed cache of the embed
    ├── embed_stub.go             # empty embed slice when -tags embed_osquery is off
    ├── embed_<os>_<arch>.go      # per-platform `//go:embed` of osqueryd
    ├── handlers.go               # endpoint list + handlers
    └── handlers_test.go          # unit tests for the endpoint surface
```

End-to-end tests live under `test/e2e/osquery-plugin/` — run with
`mage e2eOsquery` (requires DART + LXD).
