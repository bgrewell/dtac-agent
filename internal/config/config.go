package config

import (
	"fmt"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InternalSettings are never written out to the configuration file
type InternalSettings struct {
	ProductName string `json:"-" yaml:"-" mapstructure:"product_name"`
	ShortName   string `json:"-" yaml:"-" mapstructure:"short_name"`
	FileName    string `json:"-" yaml:"-" mapstructure:"file_name"`
}

// BlockingEntry is the struct for a blocking entry
type BlockingEntry struct {
	Trigger       string `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
	Detect        string `json:"detect" yaml:"detect" mapstructure:"detect"`
	TimeoutMs     int    `json:"timeout_ms" yaml:"timeout_ms" mapstructure:"timeout_ms"`
	TimeoutAction string `json:"timeout_action" yaml:"timeout_action" mapstructure:"timeout_action"`
}

// TLSConfigurationEntry is the struct for a TLS configuration
type TLSConfigurationEntry struct {
	Enabled         bool     `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Type            string   `json:"type" yaml:"type" mapstructure:"type"`
	CreateIfMissing bool     `json:"create_if_missing" yaml:"create_if_missing" mapstructure:"create_if_missing"`
	Domains         []string `json:"domains" yaml:"domains" mapstructure:"domains"`
	CertFile        string   `json:"cert" yaml:"cert" mapstructure:"cert"`
	KeyFile         string   `json:"key" yaml:"key" mapstructure:"key"`
	CAFile          string   `json:"ca" yaml:"ca" mapstructure:"ca"`
}

// LockoutEntry is the struct for a lockout entry
type LockoutEntry struct {
	Enabled        bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	AutoUnlockTime string `json:"auto_unlock_time" yaml:"auto_unlock_time" mapstructure:"auto_unlock_time"`
}

// PluginEntry is the struct for a plugin entry
type PluginEntry struct {
	Enabled          bool                             `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	PluginDir        string                           `json:"dir" yaml:"dir" mapstructure:"dir"`
	PluginGroup      string                           `json:"group" yaml:"group" mapstructure:"group"`
	LoadUnconfigured bool                             `json:"load_unconfigured" yaml:"load_unconfigured" mapstructure:"load_unconfigured"`
	TLS              TLSSelection                     `json:"tls" yaml:"tls" mapstructure:"tls"`
	Entries          map[string]*plugins.PluginConfig `json:"entries" yaml:"entries" mapstructure:"entries"`
}

// RouteEntry is the struct for a route entry
type RouteEntry struct {
	Name     string         `json:"name" yaml:"name" mapstructure:"name"`
	Route    string         `json:"route" yaml:"route" mapstructure:"route"`
	Source   *SourceEntry   `json:"source" yaml:"source" mapstructure:"source"`
	Methods  []string       `json:"methods" yaml:"methods" mapstructure:"methods"`
	Access   string         `json:"access" yaml:"access" mapstructure:"access"`
	Blocking *BlockingEntry `json:"blocking" yaml:"blocking" mapstructure:"blocking"`
	WrapJSON bool           `json:"wrap_json" yaml:"wrap_json" mapstructure:"wrap_json"`
	Mode     string         `json:"mode" yaml:"mode" mapstructure:"mode"`
}

// SourceEntry is the struct for a source entry
type SourceEntry struct {
	Type  string `json:"type" yaml:"type" mapstructure:"type"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
	RunAs string `json:"run_as" yaml:"run_as" mapstructure:"run_as"`
}

// SubsystemEntry is the struct for a subsystem entry
type SubsystemEntry struct {
	Auth       bool `json:"auth" yaml:"auth" mapstructure:"auth"`
	Diag       bool `json:"diag" yaml:"diag" mapstructure:"diag"`
	Echo       bool `json:"echo" yaml:"echo" mapstructure:"echo"`
	Hardware   bool `json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	Network    bool `json:"network" yaml:"network" mapstructure:"network"`
	Validation bool `json:"validation" yaml:"validation" mapstructure:"validation"`
}

// UpdaterEntry is the struct for an updater entry
type UpdaterEntry struct {
	Enabled         bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Token           string `json:"token" yaml:"token" mapstructure:"token"`
	Mode            string `json:"mode" yaml:"mode" mapstructure:"mode"`
	Interval        string `json:"interval" yaml:"interval" mapstructure:"interval"`
	ErrorFallback   string `json:"error_fallback" yaml:"error_fallback" mapstructure:"error_fallback"`
	RestartOnUpdate bool   `json:"restart_on_update" yaml:"restart_on_update" mapstructure:"restart_on_update"`
}

// WatchdogEntry is the struct for a watchdog entry
type WatchdogEntry struct {
	Enabled      bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	PollInterval string `json:"poll_interval" yaml:"poll_interval" mapstructure:"poll_interval"`
	Profile      string `json:"profile" yaml:"profile" mapstructure:"profile"`
	BSSID        string `json:"bssid" yaml:"bssid" mapstructure:"bssid"`
}

// AuthEntry is the struct for an auth entry
type AuthEntry struct {
	User          string `json:"admin" yaml:"admin" mapstructure:"admin"`
	Pass          string `json:"pass" yaml:"pass" mapstructure:"pass"`
	DefaultSecure bool   `json:"default_secure" yaml:"default_secure" mapstructure:"default_secure"`
	Model         string `json:"model" yaml:"model" mapstructure:"model"`
	Policy        string `json:"policy" yaml:"policy" mapstructure:"policy"`
}

// OutputEntry is the struct for an output entry
type OutputEntry struct {
	LogLevel       string `json:"log_level" yaml:"log_level" mapstructure:"log_level"`
	WrapResponses  bool   `json:"wrap_responses" yaml:"wrap_responses" mapstructure:"wrap_responses"`
	IncludeSchemas bool   `json:"include_schemas" yaml:"include_schemas" mapstructure:"include_schemas"`
}

// TLSSelection is the struct for a tls selection
type TLSSelection struct {
	Enabled bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Profile string `json:"profile" yaml:"profile" mapstructure:"profile"`
}

// CORSConfig is the struct for CORS configuration
type CORSConfig struct {
	Enabled          bool     `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins" mapstructure:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods" yaml:"allowed_methods" mapstructure:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers" yaml:"allowed_headers" mapstructure:"allowed_headers"`
	ExposedHeaders   []string `json:"exposed_headers" yaml:"exposed_headers" mapstructure:"exposed_headers"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials" mapstructure:"allow_credentials"`
	MaxAge           int      `json:"max_age" yaml:"max_age" mapstructure:"max_age"`
}

// APIEntries is the struct for a api entries
type APIEntries struct {
	REST RESTAPIEntry `json:"rest" yaml:"rest" mapstructure:"rest"`
	GRPC GRPCAPIEntry `json:"grpc" yaml:"grpc" mapstructure:"grpc"`
	JSON JSONAPIEntry `json:"json" yaml:"json" mapstructure:"json"`
}

// RESTAPIEntry is the struct for an api entry
type RESTAPIEntry struct {
	Enabled bool         `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Port    int          `json:"port" yaml:"port" mapstructure:"port"`
	TLS     TLSSelection `json:"tls" yaml:"tls" mapstructure:"tls"`
	CORS    CORSConfig   `json:"cors" yaml:"cors" mapstructure:"cors"`
}

// JSONAPIEntry is the struct for an api entry
type JSONAPIEntry struct {
	Enabled bool         `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Port    int          `json:"port" yaml:"port" mapstructure:"port"`
	TLS     TLSSelection `json:"tls" yaml:"tls" mapstructure:"tls"`
}

// GRPCAPIEntry is the struct for an api entry
type GRPCAPIEntry struct {
	Enabled    bool         `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Port       int          `json:"port" yaml:"port" mapstructure:"port"`
	Reflection bool         `json:"reflection" yaml:"reflection" mapstructure:"reflection"`
	TLS        TLSSelection `json:"tls" yaml:"tls" mapstructure:"tls"`
}

// Configuration is the struct for the configuration
type Configuration struct {
	Include         []string                         `json:"include" yaml:"include" mapstructure:"include"`
	APIs            APIEntries                       `json:"apis" yaml:"apis" mapstructure:"apis"`
	Auth            AuthEntry                        `json:"auth" yaml:"auth" mapstructure:"auth"`
	Internal        InternalSettings                 `json:"-" yaml:"-" mapstructure:"internal"`
	Lockout         LockoutEntry                     `json:"lockout" yaml:"lockout" mapstructure:"lockout"`
	Subsystems      SubsystemEntry                   `json:"subsystems" yaml:"subsystems" mapstructure:"subsystems"`
	TLS             map[string]TLSConfigurationEntry `json:"tls" yaml:"tls" mapstructure:"tls"`
	Updater         UpdaterEntry                     `json:"updater" yaml:"updater" mapstructure:"updater"`
	WifiWatchdog    WatchdogEntry                    `json:"wifi_watchdog" yaml:"wifi_watchdog" mapstructure:"wifi_watchdog"`
	Plugins         PluginEntry                      `json:"plugins" yaml:"plugins" mapstructure:"plugins"`
	CustomEndpoints []map[string]*RouteEntry         `json:"custom_endpoints" yaml:"custom_endpoints" mapstructure:"custom_endpoints"` //TODO: Needs to be updated for new architecture
	Output          OutputEntry                      `json:"output" yaml:"output" mapstructure:"output"`
	logger          *zap.Logger
}

func NewConfiguration(log *zap.Logger) (*Configuration, error) {
	// Setup viper (existing)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(GlobalConfigLocation)
	viper.AddConfigPath(LocalConfigLocation)
	viper.AddConfigPath(".")
	viper.SetConfigPermissions(0600)

	// Defaults (existing)
	for k, v := range DefaultConfig() {
		viper.SetDefault(k, v)
	}

	// Read base config (existing create-if-missing logic unchanged) ...
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read configuration file: %v", err)
		}
		// your existing create-default logic...
		log.Warn("configuration file not found")
		if ensureDir(GlobalConfigLocation, true) && checkWriteAccess(GlobalConfigLocation) {
			log.Info("[creating configuration file", zap.String("location", GlobalConfigLocation))
			if err := os.MkdirAll(GlobalConfigLocation, 0700); err != nil {
				log.Error("error creating global config directory", zap.Error(err))
			}
			if err := writeConfigWithoutInternalKeys(path.Join(GlobalConfigLocation, "config.yaml")); err != nil {
				return nil, fmt.Errorf("failed to write log file: %v", err)
			}
		} else {
			location := strings.Replace(LocalConfigLocation, "$HOME", os.Getenv("HOME"), 1)
			if err := os.MkdirAll(location, 0700); err != nil {
				log.Error("error creating user config directory", zap.Error(err))
			}
			log.Info("creating configuration file", zap.String("filename", path.Join(location, "config.yaml")))
			if err := writeConfigWithoutInternalKeys(path.Join(location, "config.yaml")); err != nil {
				return nil, fmt.Errorf("failed to write log file: %v", err)
			}
		}
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read configuration file: %v", err)
		}
	}

	// === NEW: merge include fragments (SSH-like: last wins) ===
	merged := viper.AllSettings()

	// Pull include list (from the already-read base settings or defaults)
	includes := viper.GetStringSlice("include")
	for _, pattern := range includes {
		// Expand env + ~
		expanded := os.ExpandEnv(expandUser(pattern))

		// Glob (ignore non-matches)
		matches, _ := filepath.Glob(expanded)
		if len(matches) == 0 {
			continue
		}

		// Deterministic order (like ssh: lexicographic)
		sort.Strings(matches)

		// For each file: read via a fresh viper and deep-merge into merged
		for _, f := range matches {
			vf := viper.New()
			vf.SetConfigFile(f)
			if err := vf.ReadInConfig(); err != nil {
				log.Warn("failed to read include fragment", zap.String("file", f), zap.Error(err))
				continue
			}
			merged = deepMergeMaps(merged, vf.AllSettings())
		}
	}

	// Unmarshal once into final struct
	v := viper.New()
	for k, val := range merged {
		v.Set(k, val)
	}
	var c Configuration
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	// Logger + route registration (existing)
	c.logger = log
	c.register()

	// Optional: keep watching the main file (fragments won’t auto-reload unless you add a custom watcher)
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("config file changed", zap.String("filename", e.Name))
	})
	viper.WatchConfig()

	return &c, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() map[string]interface{} {
	// Get the hostname and domain
	hostname, _ := os.Hostname()

	return map[string]interface{}{
		"include": []string{
			"$XDG_CONFIG_HOME/dtac/config.d/*.yaml",
			"~/.config/dtac/config.d/*.yaml",
			"/etc/dtac/config.d/*.yaml",
		},
		"auth.admin":                    "admin",
		"auth.pass":                     "need_to_generate_a_random_password_on_install_or_first_run",
		"auth.default_secure":           true,
		"auth.model":                    DefaultAuthModelName,
		"auth.policy":                   DefaultAuthPolicyName,
		"internal.product_name":         "DTAC Agent",
		"internal.short_name":           "dtac",
		"internal.file_name":            "dtac-agentd",
		"apis.rest.enabled":             true,
		"apis.rest.port":                8180,
		"apis.rest.tls.enabled":         true,
		"apis.rest.tls.profile":         "default",
		"apis.rest.cors.enabled":        false,
		"apis.rest.cors.allowed_origins": []string{"*"},
		"apis.rest.cors.allowed_methods": []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		"apis.rest.cors.allowed_headers": []string{"Origin", "Content-Type", "Accept", "Authorization"},
		"apis.rest.cors.exposed_headers": []string{"Content-Length"},
		"apis.rest.cors.allow_credentials": false,
		"apis.rest.cors.max_age":        3600,
		"apis.grpc.enabled":             true,
		"apis.grpc.port":                8181,
		"apis.grpc.reflection":          false,
		"apis.grpc.tls.enabled":         true,
		"apis.grpc.tls.profile":         "default",
		"apis.json.enabled":             false,
		"apis.json.port":                8182,
		"apis.json.tls.enabled":         true,
		"apis.json.tls.profile":         "default",
		"tls.default.enabled":           true,
		"tls.default.type":              "self-signed",
		"tls.default.ca":                DefaultTLSCACertName,
		"tls.default.cert":              DefaultTLSCertName,
		"tls.default.key":               DefaultTLSKeyName,
		"tls.default.create_if_missing": true,
		"tls.default.domains":           []string{"localhost", hostname},
		"lockout.enabled":               true,
		"lockout.auto_unlock_time":      "10s",
		"wifi_watchdog.enabled":         false,
		"wifi_watchdog.poll_interval":   "10s",
		"wifi_watchdog.profile":         "",
		"wifi_watchdog.bssid":           "",
		"updater.enabled":               false,
		"updater.token":                 "",
		"updater.mode":                  "auto",
		"updater.interval":              "1m",
		"updater.error_fallback":        "1h",
		"updater.restart_on_update":     true,
		"plugins.enabled":               true,
		"plugins.dir":                   DefaultPluginLocation,
		"plugins.group":                 "plugins",
		"plugins.load_unconfigured":     false,
		"plugins.tls.enabled":           true,
		"plugins.tls.profile":           "default",
		"subsystems.auth":               true,
		"subsystems.diag":               true,
		"subsystems.echo":               true,
		"subsystems.network":            true,
		"subsystems.hardware":           true,
		"subsystems.validation":         true,
		"custom_endpoints":              []map[string]*RouteEntry{},
		"output.log_level":              "debug",
		"output.include_schemas":        false,
		"output.wrap_responses":         false,
	}
}

// register registers the routes that this module handles
func (c *Configuration) register() {
	//// Create a group for this subsystem
	//base := c.router.Group("config")
	//
	//// Routes
	//routes := []types.RouteInfo{
	//	{HTTPMethod: http.MethodGet, Path: "/", Handler: c.rootHandler},
	//}
	//
	//// Register routes
	//for _, route := range routes {
	//	base.Handle(route.HTTPMethod, route.Path, route.Handler)
	//}
	//c.logger.Info("registered routes", zap.Int("routes", len(routes)))
	_ = c.rootHandler
}

func (c *Configuration) rootHandler(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"configuration": c,
	})
}

// ensureDir checks if the directory exists. If it doesn't and the `create` flag is true, it attempts to create it.
// It returns true if the directory exists or was created successfully, and false otherwise.
func ensureDir(dir string, create bool) bool {
	if _, err := os.Stat(dir); err == nil {
		// Directory exists
		return true
	} else if os.IsNotExist(err) {
		// Directory doesn't exist. Check the `create` flag.
		if create {
			// Try to create the directory
			if mkdirErr := os.MkdirAll(dir, 0755); mkdirErr != nil {
				return false
			}
			return true
		}
		return false
	} else {
		// Some other error occurred
		return false
	}
}

// CheckWriteAccess checks if we have write access to a directory.
func checkWriteAccess(dir string) bool {
	// Try to create a temporary file in the directory.
	tempFile, err := os.CreateTemp(dir, "write-check-")
	if err != nil {
		return false
	}

	// If successful, delete the temporary file and return true.
	defer os.Remove(tempFile.Name())
	return true
}

func writeConfigWithoutInternalKeys(filename string) error {
	// Backup internal settings
	internalBackup := make(map[string]interface{})
	allSettings := viper.AllSettings()
	for key, value := range allSettings {
		if key == "internal" {
			internalBackup[key] = value
			viper.SetDefault(key, nil) // Temporarily remove the settings
		}
	}

	// Add default plugin entries
	pluginKVPs := map[string]interface{}{
		"plugins.entries.hello.enabled":             true,
		"plugins.entries.hello.hash":                "",
		"plugins.entries.hello.user":                "",
		"plugins.entries.hello.config.message":      "hello world plugin",
		"plugins.entries.iperf.enabled":             true,
		"plugins.entries.iperf.hash":                "",
		"plugins.entries.iperf.user":                "",
		"plugins.entries.iperf.config.bind":         "0.0.0.0",
		"plugins.entries.iperf.config.control_port": 8191,
	}
	for key, value := range pluginKVPs {
		viper.SetDefault(key, value)
	}

	// Write the configuration without internal keys
	err := viper.WriteConfigAs(filename)
	if err != nil {
		return err
	}

	// Restore internal settings from the backup
	for key, value := range internalBackup {
		viper.SetDefault(key, value)
	}

	return nil
}

// GetConfigValue gets a value from the configuration.
func GetConfigValue(cfg *Configuration, key string) (interface{}, error) {
	keys := strings.Split(key, ".")
	var current interface{} = cfg

	for _, k := range keys {
		value, found := findField(current, k)
		if !found {
			return nil, fmt.Errorf("key not found: %s", key)
		}

		current = value
	}

	return current, nil
}

// findField finds a field in a struct by its JSON or YAML tag.
func findField(v interface{}, key string) (interface{}, bool) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Kind() == reflect.Struct {
		// Look up the field using the struct tags
		field := value.FieldByNameFunc(func(fieldName string) bool {
			field, _ := value.Type().FieldByName(fieldName)
			tag := field.Tag
			jsonTag := tag.Get("json")
			yamlTag := tag.Get("yaml")
			return jsonTag == key || yamlTag == key
		})

		if field.IsValid() {
			return field.Interface(), true
		}
	}

	return nil, false
}

// expandUser expands a leading "~" to the current user's home directory.
func expandUser(p string) string {
	if strings.HasPrefix(p, "~") {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			return filepath.Join(home, strings.TrimPrefix(p, "~"))
		}
	}
	return p
}

// deepMergeMaps merges b into a (recursively). Later values override earlier (“last wins”).
func deepMergeMaps(a, b map[string]any) map[string]any {
	out := make(map[string]any, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, bv := range b {
		if av, ok := out[k]; ok {
			am, aok := av.(map[string]any)
			bm, bok := bv.(map[string]any)
			if aok && bok {
				out[k] = deepMergeMaps(am, bm)
				continue
			}
		}
		out[k] = bv
	}
	return out
}
