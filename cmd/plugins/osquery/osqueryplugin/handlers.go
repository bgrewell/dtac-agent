package osqueryplugin

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// cannedQuery declares one of the curated read-only endpoints whose SQL does
// not depend on the host OS. The handler is generated from `sql` verbatim.
type cannedQuery struct {
	path        string
	description string
	authGroup   endpoint.AuthGroup
	sql         string
}

// dispatchedQuery declares a curated endpoint whose SQL is selected at
// handler time based on runtime.GOOS — used for things like `services` or
// `installed_packages` where the relevant osquery table differs per platform.
// Platforms that are absent from `sqlByOS` return an explicit error.
type dispatchedQuery struct {
	path        string
	description string
	authGroup   endpoint.AuthGroup
	sqlByOS     map[string]string
}

// curatedQueries is the small set of pre-baked read-only endpoints shipped
// with the plugin. Adding here is the cheap, high-leverage way to expose new
// host signals — anything more bespoke (parameterized, expensive, mutating)
// belongs behind /query with stricter auth.
//
// Auth group rationale: User-tier for identity-only facts that any caller
// would want; Operator-tier for inventory/topology; Admin-tier for audit
// surfaces or anything that leaks credentials/scheduling.
var curatedQueries = []cannedQuery{
	// ─── original set ──────────────────────────────────────────────────────
	{
		path:        "processes",
		description: "running processes (pid, name, path, cmdline, uid, parent)",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT pid, name, path, cmdline, uid, parent FROM processes",
	},
	{
		path:        "users",
		description: "local user accounts",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM users",
	},
	{
		path:        "listening_ports",
		description: "sockets the host is listening on",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM listening_ports",
	},
	{
		path:        "os_version",
		description: "operating system version",
		authGroup:   endpoint.AuthGroupUser,
		sql:         "SELECT * FROM os_version",
	},
	{
		path:        "interface_addresses",
		description: "addresses configured on network interfaces",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM interface_addresses",
	},
	{
		path:        "logged_in_users",
		description: "utmp-backed login sessions — TTYs only; for transit/forward-only SSH use /active_sessions",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM logged_in_users",
	},
	// ─── tier 1 ────────────────────────────────────────────────────────────
	{
		path:        "system_info",
		description: "host identity — hostname, hardware vendor/model, CPU, RAM",
		authGroup:   endpoint.AuthGroupUser,
		sql:         "SELECT * FROM system_info",
	},
	{
		path:        "kernel_info",
		description: "kernel version, architecture, and build path",
		authGroup:   endpoint.AuthGroupUser,
		sql:         "SELECT * FROM kernel_info",
	},
	{
		path:        "uptime",
		description: "system uptime in days, hours, minutes, seconds",
		authGroup:   endpoint.AuthGroupUser,
		sql:         "SELECT * FROM uptime",
	},
	{
		path:        "routes",
		description: "IPv4/IPv6 routing table",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM routes",
	},
	{
		path:        "kernel_modules",
		description: "loaded kernel modules (Linux only; returns empty on darwin)",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM kernel_modules",
	},
	{
		path:        "crontab",
		description: "scheduled cron jobs across all user crontabs (root needed for non-self entries)",
		authGroup:   endpoint.AuthGroupAdmin,
		sql:         "SELECT * FROM crontab",
	},
	{
		path:        "process_open_sockets",
		description: "process → open socket mapping; full visibility requires root",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM process_open_sockets",
	},
	{
		// SSH-aware session enumeration. Each authenticated SSH connection
		// (including TTY-less command exec, -N port-forward only, ProxyJump
		// transits, scp/sftp, and multiplexed slaves) shows up as an sshd
		// child process with an established remote socket. utmp-based tools
		// (who/w/logged_in_users) miss all of these except interactive
		// shells; this query catches them all. Cross-platform — works on
		// Linux and darwin without changes.
		path:        "active_sessions",
		description: "all active SSH sessions including TTY-less and transit (jump-host) connections; needs root for non-self users",
		authGroup:   endpoint.AuthGroupOperator,
		sql: `SELECT
				p.pid           AS sshd_pid,
				p.parent        AS sshd_parent,
				p.uid           AS uid,
				u.username      AS username,
				s.remote_address AS remote_addr,
				s.remote_port    AS remote_port,
				s.local_address  AS local_addr,
				s.local_port     AS local_port,
				p.start_time    AS start_time
			FROM processes p
			LEFT JOIN process_open_sockets s
				ON p.pid = s.pid AND s.family IN (2, 10) AND s.remote_address != ''
			LEFT JOIN users u ON p.uid = u.uid
			WHERE p.name = 'sshd'
			  AND p.uid != 0
			  AND s.remote_address IS NOT NULL
			ORDER BY p.uid, p.start_time`,
	},
	// ─── tier 2 ────────────────────────────────────────────────────────────
	{
		path:        "groups",
		description: "local groups",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM groups",
	},
	{
		path:        "last_logins",
		description: "wtmp-backed historical logins (audit-relevant; subject to wtmp rotation)",
		authGroup:   endpoint.AuthGroupAdmin,
		sql:         "SELECT * FROM last WHERE host != ''",
	},
	{
		path:        "arp_cache",
		description: "ARP / IPv6 neighbor cache",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM arp_cache",
	},
	{
		path:        "dns_resolvers",
		description: "configured DNS resolvers (/etc/resolv.conf and friends)",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM dns_resolvers",
	},
	{
		path:        "etc_hosts",
		description: "static host-name to IP mappings from /etc/hosts",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM etc_hosts",
	},
	{
		path:        "block_devices",
		description: "block-device inventory (Linux only; returns empty on darwin)",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM block_devices",
	},
	{
		path:        "mounts",
		description: "active filesystem mounts",
		authGroup:   endpoint.AuthGroupOperator,
		sql:         "SELECT * FROM mounts",
	},
	{
		path:        "sudoers",
		description: "parsed /etc/sudoers and drop-ins (requires root to read sudoers files)",
		authGroup:   endpoint.AuthGroupAdmin,
		sql:         "SELECT * FROM sudoers",
	},
	{
		path:        "osquery_info",
		description: "bundled osqueryd self-status (version, uptime, watcher state)",
		authGroup:   endpoint.AuthGroupUser,
		sql:         "SELECT * FROM osquery_info",
	},
}

// dispatchedQueries fan out one logical endpoint across multiple OS-specific
// osquery tables. Returning a clear error on unsupported platforms beats
// silently returning empty rows — the caller learns the endpoint is not
// applicable rather than concluding the host has zero services / packages.
var dispatchedQueries = []dispatchedQuery{
	{
		path:        "services",
		description: "active system services — systemd_units on Linux, launchd on darwin",
		authGroup:   endpoint.AuthGroupOperator,
		sqlByOS: map[string]string{
			"linux": `SELECT id, description, load_state, active_state, sub_state, fragment_path, user
			           FROM systemd_units
			           WHERE active_state = 'active'`,
			"darwin": `SELECT label, program, program_arguments, run_at_load, keep_alive, disabled
			            FROM launchd`,
		},
	},
	{
		path:        "installed_packages",
		description: "installed software packages — deb+rpm on Linux, apps on darwin",
		authGroup:   endpoint.AuthGroupOperator,
		sqlByOS: map[string]string{
			"linux": `SELECT name, version, arch, source, 'deb' AS package_format FROM deb_packages
			           UNION ALL
			           SELECT name, version, arch, source, 'rpm' AS package_format FROM rpm_packages`,
			"darwin": `SELECT name, bundle_version AS version, bundle_short_version AS short_version,
			                  bundle_identifier AS bundle_id, path
			            FROM apps`,
		},
	},
}

// queryRequest is the JSON body shape POSTed to the generic /query endpoint.
type queryRequest struct {
	SQL string `json:"sql"`
}

// buildEndpoints assembles the full endpoint slice that Register reports
// back to the agent. It is broken out from Register so it can be exercised
// from tests without a live osqueryd.
func (h *OsqueryPlugin) buildEndpoints(defaultSecure bool) []*endpoint.Endpoint {
	endpoints := make([]*endpoint.Endpoint, 0, len(curatedQueries)+len(dispatchedQueries)+1)

	endpoints = append(endpoints, endpoint.NewEndpoint(
		"query",
		endpoint.ActionCreate,
		"run an arbitrary osquery SQL statement and return rows",
		h.handleQuery,
		defaultSecure,
		endpoint.AuthGroupAdmin.String(),
		endpoint.WithBody(&queryRequest{}),
		endpoint.WithOutput(&queryEnvelope{}),
	))

	for _, q := range curatedQueries {
		q := q
		endpoints = append(endpoints, endpoint.NewEndpoint(
			q.path,
			endpoint.ActionRead,
			q.description,
			func(in *endpoint.Request) (*endpoint.Response, error) {
				return h.runQueryWrapped(in, q.sql)
			},
			defaultSecure,
			q.authGroup.String(),
			endpoint.WithOutput(&queryEnvelope{}),
		))
	}

	for _, d := range dispatchedQueries {
		d := d
		endpoints = append(endpoints, endpoint.NewEndpoint(
			d.path,
			endpoint.ActionRead,
			d.description,
			func(in *endpoint.Request) (*endpoint.Response, error) {
				sql, ok := d.sqlByOS[runtime.GOOS]
				if !ok {
					return nil, fmt.Errorf("endpoint /%s is not implemented for %s", d.path, runtime.GOOS)
				}
				return h.runQueryWrapped(in, sql)
			},
			defaultSecure,
			d.authGroup.String(),
			endpoint.WithOutput(&queryEnvelope{}),
		))
	}

	return endpoints
}

// handleQuery serves POST /query. The body must be `{"sql":"..."}`.
func (h *OsqueryPlugin) handleQuery(in *endpoint.Request) (*endpoint.Response, error) {
	if len(in.Body) == 0 {
		return nil, fmt.Errorf("empty request body; expected {\"sql\":\"...\"}")
	}
	var req queryRequest
	if err := json.Unmarshal(in.Body, &req); err != nil {
		return nil, fmt.Errorf("decoding query request: %w", err)
	}
	if req.SQL == "" {
		return nil, fmt.Errorf("query request missing required field \"sql\"")
	}
	return h.runQueryWrapped(in, req.SQL)
}
