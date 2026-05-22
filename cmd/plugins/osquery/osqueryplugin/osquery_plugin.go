package osqueryplugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
	"github.com/osquery/osquery-go"
)

// Compile-time assertion that OsqueryPlugin satisfies the plugin contract.
var _ plugins.Plugin = &OsqueryPlugin{}

// OsqueryPlugin wraps a private osqueryd subprocess and exposes its query
// surface through the dtac plugin framework.
type OsqueryPlugin struct {
	plugins.PluginBase

	mu       sync.Mutex
	cfg      Config
	daemon   *daemon
	client   *osquery.ExtensionManagerClient
	maxRows  int
	timeout  time.Duration
	stopOnce sync.Once
}

// NewOsqueryPlugin returns a fresh, un-started plugin instance. The osqueryd
// subprocess is not spawned here — Register() is the entry point for that so
// that the agent's JSON config can be honored.
func NewOsqueryPlugin() *OsqueryPlugin {
	p := &OsqueryPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}
	p.SetRootPath("osquery")
	return p
}

// Name returns the plugin's type name. It must be defined on OsqueryPlugin
// (not promoted from PluginBase) so the agent sees the right value during
// CONNECT handshake.
func (h *OsqueryPlugin) Name() string {
	return "OsqueryPlugin"
}

// Register parses the JSON config, spawns osqueryd, opens an osquery-go
// client, and publishes the plugin's endpoints back to the agent.
func (h *OsqueryPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	cfg, err := parseConfig(request.Config)
	if err != nil {
		return fmt.Errorf("parsing osquery plugin config: %w", err)
	}
	h.cfg = cfg
	h.maxRows = cfg.MaxRows
	h.timeout = cfg.QueryTimeout

	stdoutSink := func(line string) {
		h.Log(plugins.LevelInfo, line, map[string]string{"source": "osqueryd", "stream": "stdout"})
	}
	stderrSink := func(line string) {
		h.Log(plugins.LevelWarning, line, map[string]string{"source": "osqueryd", "stream": "stderr"})
	}

	d, err := startDaemon(cfg, stdoutSink, stderrSink)
	if err != nil {
		return fmt.Errorf("starting osqueryd: %w", err)
	}
	h.daemon = d

	client, err := osquery.NewClient(d.SocketPath(), cfg.StartupTimeout)
	if err != nil {
		_ = d.Stop(context.Background())
		return fmt.Errorf("connecting to osqueryd: %w", err)
	}
	h.client = client

	endpoints := h.buildEndpoints(request.DefaultSecure)
	h.RegisterMethods(endpoints)
	for _, ep := range endpoints {
		reply.Endpoints = append(reply.Endpoints, utility.ConvertEndpointToPluginEndpoint(ep))
	}

	h.Log(plugins.LevelInfo, "osquery plugin registered", map[string]string{
		"endpoint_count": strconv.Itoa(len(endpoints)),
		"socket":         d.SocketPath(),
	})
	return nil
}

// Shutdown tears down the osqueryd subprocess. The plugin framework has no
// Stop() hook yet (see issue #17) so the entrypoint's signal handler must
// call this explicitly.
func (h *OsqueryPlugin) Shutdown(ctx context.Context) {
	h.stopOnce.Do(func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		if h.client != nil {
			h.client.Close()
			h.client = nil
		}
		if h.daemon != nil {
			_ = h.daemon.Stop(ctx)
			h.daemon = nil
		}
	})
}

// queryRows runs sql against osqueryd with the plugin's configured timeout
// and row cap, returning the rows that the JSON marshaler can hand back to
// the agent.
func (h *OsqueryPlugin) queryRows(sql string) ([]map[string]string, error) {
	h.mu.Lock()
	client := h.client
	h.mu.Unlock()
	if client == nil {
		return nil, fmt.Errorf("osquery client is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()
	rows, err := client.QueryRowsContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	if h.maxRows > 0 && len(rows) > h.maxRows {
		rows = rows[:h.maxRows]
	}
	return rows, nil
}

// queryEnvelope is the shape returned by both /query and the curated
// endpoints, so clients can rely on a single response format.
type queryEnvelope struct {
	Rows       []map[string]string `json:"rows"`
	RowCount   int                 `json:"row_count"`
	DurationMs int64               `json:"duration_ms"`
	Truncated  bool                `json:"truncated,omitempty"`
}

// runQueryWrapped is the shared handler core for both /query and the curated
// endpoints. It does the timing, row-cap, and JSON marshaling so callers can
// stay one-liners.
func (h *OsqueryPlugin) runQueryWrapped(in *endpoint.Request, sql string) (*endpoint.Response, error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		start := time.Now()
		rows, err := h.queryRows(sql)
		if err != nil {
			return nil, nil, err
		}
		env := queryEnvelope{
			Rows:       rows,
			RowCount:   len(rows),
			DurationMs: time.Since(start).Milliseconds(),
		}
		// queryRows already truncates; report whether it did.
		if h.maxRows > 0 && len(rows) == h.maxRows {
			env.Truncated = true
		}
		body, err := json.Marshal(env)
		if err != nil {
			return nil, nil, err
		}
		headers := map[string][]string{
			"Content-Type":  {"application/json"},
			"X-PLUGIN-NAME": {h.Name()},
		}
		return headers, body, nil
	}, "osquery query result")
}
