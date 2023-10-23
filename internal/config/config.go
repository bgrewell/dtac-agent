package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bgrewell/gin-plugins/loader"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InternalSettings are never written out to the configuration file
type InternalSettings struct {
	ProductName string `json:"product_name" yaml:"product_name" mapstructure:"product_name"`
	ShortName   string `json:"short_name" yaml:"short_name" mapstructure:"short_name"`
	FileName    string `json:"file_name" yaml:"file_name" mapstructure:"file_name"`
}

// BlockingEntry is the struct for a blocking entry
type BlockingEntry struct {
	Trigger       string `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
	Detect        string `json:"detect" yaml:"detect" mapstructure:"detect"`
	TimeoutMs     int    `json:"timeout_ms" yaml:"timeout_ms" mapstructure:"timeout_ms"`
	TimeoutAction string `json:"timeout_action" yaml:"timeout_action" mapstructure:"timeout_action"`
}

// ListenerEntry is the struct for a listener entry
type ListenerEntry struct {
	Port  int                `json:"port" yaml:"port" mapstructure:"port"`
	HTTPS ListenerHTTPSEntry `json:"https" yaml:"https" mapstructure:"https"`
}

// ListenerHTTPSEntry is the struct for a listener https entry
type ListenerHTTPSEntry struct {
	Enabled         bool     `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Type            string   `json:"type" yaml:"type" mapstructure:"type"`
	CreateIfMissing bool     `json:"create_if_missing" yaml:"create_if_missing" mapstructure:"create_if_missing"`
	Domains         []string `json:"domains" yaml:"domains" mapstructure:"domains"`
	CertFile        string   `json:"cert" yaml:"cert" mapstructure:"cert"`
	KeyFile         string   `json:"key" yaml:"key" mapstructure:"key"`
}

// LockoutEntry is the struct for a lockout entry
type LockoutEntry struct {
	Enabled        bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	AutoUnlockTime string `json:"auto_unlock_time" yaml:"auto_unlock_time" mapstructure:"auto_unlock_time"`
}

// PluginEntry is the struct for a plugin entry
type PluginEntry struct {
	Enabled          bool                            `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	PluginDir        string                          `json:"dir" yaml:"dir" mapstructure:"dir"`
	PluginGroup      string                          `json:"group" yaml:"group" mapstructure:"group"`
	LoadUnconfigured bool                            `json:"load_unconfigured" yaml:"load_unconfigured" mapstructure:"load_unconfigured"`
	Entries          map[string]*loader.PluginConfig `json:"entries" yaml:"entries" mapstructure:"entries"`
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
	Diag     bool `json:"diag" yaml:"diag" mapstructure:"diag"`
	Hardware bool `json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	Network  bool `json:"network" yaml:"network" mapstructure:"network"`
	//TODO: Below this line are old and should be removed
	Firewall     bool `json:"firewall" yaml:"firewall" mapstructure:"firewall"`
	Iperf        bool `json:"iperf" yaml:"iperf" mapstructure:"iperf"`
	TCPReflector bool `json:"tcp_reflector" yaml:"tcp_reflector" mapstructure:"tcp_reflector"`
	UDPReflector bool `json:"udp_reflector" yaml:"udp_reflector" mapstructure:"udp_reflector"`
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
	User          string `json:"user" yaml:"user" mapstructure:"user"`
	Pass          string `json:"pass" yaml:"pass" mapstructure:"pass"`
	DefaultSecure bool   `json:"default_secure" yaml:"default_secure" mapstructure:"default_secure"`
	Model         string `json:"model" yaml:"model" mapstructure:"model"`
	Policy        string `json:"policy" yaml:"policy" mapstructure:"policy"`
}

// Configuration is the struct for the configuration
type Configuration struct {
	Auth         AuthEntry                `json:"authn" yaml:"authn" mapstructure:"authn"`
	Internal     InternalSettings         `json:"internal" yaml:"internal" mapstructure:"internal"`
	Listener     ListenerEntry            `json:"listener" yaml:"listener" mapstructure:"listener"`
	Lockout      LockoutEntry             `json:"lockout" yaml:"lockout" mapstructure:"lockout"`
	Subsystems   SubsystemEntry           `json:"subsystems" yaml:"subsystems" mapstructure:"subsystems"`
	Updater      UpdaterEntry             `json:"updater" yaml:"updater" mapstructure:"updater"`
	WifiWatchdog WatchdogEntry            `json:"wifi_watchdog" yaml:"wifi_watchdog" mapstructure:"wifi_watchdog"`
	Plugins      PluginEntry              `json:"plugins" yaml:"plugins" mapstructure:"plugins"`
	Custom       []map[string]*RouteEntry `json:"routes" yaml:"routes" mapstructure:"routes"`
	router       *gin.Engine
	logger       *zap.Logger
}

// NewConfiguration creates a new configuration
func NewConfiguration(router *gin.Engine, log *zap.Logger) (config *Configuration, err error) {
	// Setup configuration file location(s)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(GlobalConfigLocation)
	viper.AddConfigPath(LocalConfigLocation)
	viper.AddConfigPath(".")
	viper.SetConfigPermissions(0600)

	// Get the hostname and domain
	hostname, _ := os.Hostname()

	// Setup default values
	kvp := map[string]interface{}{
		"authn.user":                           "admin",
		"authn.pass":                           "need_to_generate_a_random_password_on_install_or_first_run",
		"authn.default_secure":                 true,
		"authn.model":                          DefaultAuthModelName,
		"authn.policy":                         DefaultAuthPolicyName,
		"internal.product_name":                "DTAC Agent",
		"internal.short_name":                  "dtac",
		"internal.file_name":                   "dtac-agentd",
		"listener.port":                        8180,
		"listener.https.enabled":               true,
		"listener.https.type":                  "self-signed",
		"listener.https.create_if_missing":     true,
		"listener.https.domains":               []string{"localhost", "127.0.0.1", hostname},
		"listener.https.cert":                  DefaultTLSCertName,
		"listener.https.key":                   DefaultTLSKeyName,
		"lockout.enabled":                      true,
		"lockout.auto_unlock_time":             "10s",
		"wifi_watchdog.enabled":                false,
		"wifi_watchdog.poll_interval":          "10s",
		"wifi_watchdog.profile":                "",
		"wifi_watchdog.bssid":                  "",
		"updater.enabled":                      false,
		"updater.token":                        "",
		"updater.mode":                         "auto",
		"updater.interval":                     "1m",
		"updater.error_fallback":               "1h",
		"updater.restart_on_update":            true,
		"plugins.enabled":                      true,
		"plugins.dir":                          DefaultPluginLocation,
		"plugins.group":                        "plugins",
		"plugins.load_unconfigured":            false,
		"plugins.entries.hello.enabled":        true,
		"plugins.entries.hello.cookie":         "this_is_not_a_security_feature",
		"plugins.entries.hello.hash":           "",
		"plugins.entries.hello.user":           "",
		"plugins.entries.hello.config.message": "hello world plugin",
		"subsystems.diag":                      true,
		"subsystems.network":                   true,
		"subsystems.hardware":                  true,
		"routes":                               []map[string]*RouteEntry{},
	}
	for k, v := range kvp {
		viper.SetDefault(k, v)
	}

	if err := viper.ReadInConfig(); err != nil {
		// If the error is something other than the configuration file isn't found then
		// throw a fatal error for the user to handle.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read configuration file: %v", err)
		}

		// If the file is simply not found then create a default one and display a warning
		// to the user
		log.Warn("configuration file not found")
		if ensureDir(GlobalConfigLocation, true) && checkWriteAccess(GlobalConfigLocation) {
			log.Info("[creating configuration file", zap.String("location", GlobalConfigLocation))
			if err := os.MkdirAll(GlobalConfigLocation, 0700); err != nil {
				log.Error("error creating global config directory", zap.Error(err))
			}

			err := writeConfigWithoutInternalKeys(path.Join(GlobalConfigLocation, "config.yaml"))
			if err != nil {
				return nil, fmt.Errorf("failed to write log file: %v", err)
			}
		} else {
			location := strings.Replace(LocalConfigLocation, "$HOME", os.Getenv("HOME"), 1)
			if err := os.MkdirAll(location, 0700); err != nil {
				log.Error("error creating user config directory", zap.Error(err))
			}
			log.Info("creating configuration file", zap.String("filename", path.Join(location, "config.yaml")))
			err := writeConfigWithoutInternalKeys(path.Join(location, "config.yaml"))
			if err != nil {
				return nil, fmt.Errorf("failed to write log file: %v", err)
			}
		}
		err := viper.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to read configuration file: %v", err)
		}
	}

	var c Configuration
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %v", err)
	}

	// Setup routes
	c.router = router
	c.logger = log
	if err := c.Register(); err != nil {
		log.Error("failed to register config routes", zap.Error(err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("config file changed", zap.String("filename", e.Name))
	})
	viper.WatchConfig()
	return &c, nil
}

// Register registers the routes that this module handles
func (c *Configuration) Register() error {
	// Create a group for this subsystem
	base := c.router.Group("config")

	// Routes
	routes := []types.RouteInfo{
		{HTTPMethod: http.MethodGet, Path: "/", Handler: c.rootHandler},
	}

	// Register routes
	for _, route := range routes {
		base.Handle(route.HTTPMethod, route.Path, route.Handler)
	}
	c.logger.Info("registered routes", zap.Int("routes", len(routes)))

	return nil
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
		if strings.HasPrefix(key, "internal.") {
			internalBackup[key] = value
			viper.Set(key, nil) // Temporarily remove the setting
		}
	}

	// Write the configuration without internal keys
	err := viper.WriteConfigAs(filename)
	if err != nil {
		return err
	}

	// Restore internal settings from the backup
	for key, value := range internalBackup {
		viper.Set(key, value)
	}

	return nil
}
