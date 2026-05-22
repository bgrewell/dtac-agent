package osqueryplugin

import (
	"testing"

	"github.com/bgrewell/dtac-agent/pkg/endpoint"
)

// TestBuildEndpoints asserts that the plugin advertises the expected
// endpoint surface — one generic /query plus the curated read-only set —
// with the right HTTP action and auth-group mapping. It deliberately does
// not spawn osqueryd, so it stays cheap and CI-friendly.
func TestBuildEndpoints(t *testing.T) {
	p := NewOsqueryPlugin()
	eps := p.buildEndpoints(true)

	want := map[string]struct {
		action endpoint.Action
		auth   endpoint.AuthGroup
	}{
		// generic
		"query": {endpoint.ActionCreate, endpoint.AuthGroupAdmin},

		// original curated set
		"processes":           {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"users":               {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"listening_ports":     {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"os_version":          {endpoint.ActionRead, endpoint.AuthGroupUser},
		"interface_addresses": {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"logged_in_users":     {endpoint.ActionRead, endpoint.AuthGroupOperator},

		// tier 1
		"system_info":          {endpoint.ActionRead, endpoint.AuthGroupUser},
		"kernel_info":          {endpoint.ActionRead, endpoint.AuthGroupUser},
		"uptime":               {endpoint.ActionRead, endpoint.AuthGroupUser},
		"routes":               {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"kernel_modules":       {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"crontab":              {endpoint.ActionRead, endpoint.AuthGroupAdmin},
		"process_open_sockets": {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"active_sessions":      {endpoint.ActionRead, endpoint.AuthGroupOperator},

		// tier 2
		"groups":        {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"last_logins":   {endpoint.ActionRead, endpoint.AuthGroupAdmin},
		"arp_cache":     {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"dns_resolvers": {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"etc_hosts":     {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"block_devices": {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"mounts":        {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"sudoers":       {endpoint.ActionRead, endpoint.AuthGroupAdmin},
		"osquery_info":  {endpoint.ActionRead, endpoint.AuthGroupUser},

		// per-OS dispatched
		"services":           {endpoint.ActionRead, endpoint.AuthGroupOperator},
		"installed_packages": {endpoint.ActionRead, endpoint.AuthGroupOperator},
	}

	if len(eps) != len(want) {
		t.Fatalf("endpoint count: got %d, want %d", len(eps), len(want))
	}

	seen := make(map[string]bool, len(eps))
	for _, ep := range eps {
		exp, ok := want[ep.Path]
		if !ok {
			t.Errorf("unexpected endpoint path %q", ep.Path)
			continue
		}
		if ep.Action != exp.action {
			t.Errorf("endpoint %q: action = %q, want %q", ep.Path, ep.Action, exp.action)
		}
		if ep.AuthGroup != exp.auth.String() {
			t.Errorf("endpoint %q: auth_group = %q, want %q", ep.Path, ep.AuthGroup, exp.auth.String())
		}
		if !ep.Secure {
			t.Errorf("endpoint %q: secure flag should propagate defaultSecure=true", ep.Path)
		}
		if ep.Function == nil {
			t.Errorf("endpoint %q: missing handler function", ep.Path)
		}
		seen[ep.Path] = true
	}
	for path := range want {
		if !seen[path] {
			t.Errorf("missing endpoint %q", path)
		}
	}
}

// TestParseConfigDefaults asserts that empty config yields sensible defaults
// and that explicit overrides round-trip through JSON.
func TestParseConfigDefaults(t *testing.T) {
	cfg, err := parseConfig("")
	if err != nil {
		t.Fatalf("parseConfig(empty): %v", err)
	}
	if cfg.QueryTimeout <= 0 || cfg.StartupTimeout <= 0 || cfg.MaxRows <= 0 {
		t.Errorf("defaults not applied: %+v", cfg)
	}

	cfg, err = parseConfig(`{"max_rows": 5, "query_timeout": 1000000000}`)
	if err != nil {
		t.Fatalf("parseConfig(json): %v", err)
	}
	if cfg.MaxRows != 5 {
		t.Errorf("MaxRows: got %d, want 5", cfg.MaxRows)
	}
	if cfg.QueryTimeout.Seconds() != 1.0 {
		t.Errorf("QueryTimeout: got %v, want 1s", cfg.QueryTimeout)
	}
}
