package configuration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/bgrewell/gin-plugins/loader"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Config *Configuration
)

type BlockingEntry struct {
	Trigger       string `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
	Detect        string `json:"detect" yaml:"detect" mapstructure:"detect"`
	TimeoutMs     int    `json:"timeout_ms" yaml:"timeout_ms" mapstructure:"timeout_ms"`
	TimeoutAction string `json:"timeout_action" yaml:"timeout_action" mapstructure:"timeout_action"`
}

type ListenerEntry struct {
	Port  int                `json:"port" yaml:"port" mapstructure:"port"`
	Https ListenerHttpsEntry `json:"https" yaml:"https" mapstructure:"https"`
}

type ListenerHttpsEntry struct {
	Enabled         bool     `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Type            string   `json:"type" yaml:"type" mapstructure:"type"`
	CreateIfMissing bool     `json:"create_if_missing" yaml:"create_if_missing" mapstructure:"create_if_missing"`
	Domains         []string `json:"domains" yaml:"domains" mapstructure:"domains"`
	CertFile        string   `json:"cert" yaml:"cert" mapstructure:"cert"`
	KeyFile         string   `json:"key" yaml:"key" mapstructure:"key"`
}

type LockoutEntry struct {
	Enabled        bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	AutoUnlockTime string `json:"auto_unlock_time" yaml:"auto_unlock_time" mapstructure:"auto_unlock_time"`
}

type PluginEntry struct {
	Enabled   bool                            `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	PluginDir string                          `json:"dir" yaml:"dir" mapstructure:"dir"`
	Entries   map[string]*loader.PluginConfig `json:"entries" yaml:"entries" mapstructure:"entries"`
}

type RouteEntry struct {
	Name     string         `json:"name" yaml:"name" mapstructure:"name"`
	Route    string         `json:"route" yaml:"route" mapstructure:"route"`
	Source   *SourceEntry   `json:"source" yaml:"source" mapstructure:"source"`
	Methods  []string       `json:"methods" yaml:"methods" mapstructure:"methods"`
	Access   string         `json:"access" yaml:"access" mapstructure:"access"`
	Blocking *BlockingEntry `json:"blocking" yaml:"blocking" mapstructure:"blocking"`
	WrapJson bool           `json:"wrap_json" yaml:"wrap_json" mapstructure:"wrap_json"`
	Mode     string         `json:"mode" yaml:"mode" mapstructure:"mode"`
}

type SourceEntry struct {
	Type  string `json:"type" yaml:"type" mapstructure:"type"`
	Value string `json:"value" yaml:"value" mapstructure:"value"`
	RunAs string `json:"run_as" yaml:"run_as" mapstructure:"run_as"`
}

type SubsystemEntry struct {
	Firewall     bool `json:"firewall" yaml:"firewall" mapstructure:"firewall"`
	Iperf        bool `json:"iperf" yaml:"iperf" mapstructure:"iperf"`
	TcpReflector bool `json:"tcp_reflector" yaml:"tcp_reflector" mapstructure:"tcp_reflector"`
	UdpReflector bool `json:"udp_reflector" yaml:"udp_reflector" mapstructure:"udp_reflector"`
}

type UpdaterEntry struct {
	Enabled         bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Token           string `json:"token" yaml:"token" mapstructure:"token"`
	Mode            string `json:"mode" yaml:"mode" mapstructure:"mode"`
	Interval        string `json:"interval" yaml:"interval" mapstructure:"interval"`
	ErrorFallback   string `json:"error_fallback" yaml:"error_fallback" mapstructure:"error_fallback"`
	RestartOnUpdate bool   `json:"restart_on_update" yaml:"restart_on_update" mapstructure:"restart_on_update"`
}

type WatchdogEntry struct {
	Enabled      bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	PollInterval string `json:"poll_interval" yaml:"poll_interval" mapstructure:"poll_interval"`
	Profile      string `json:"profile" yaml:"profile" mapstructure:"profile"`
	BSSID        string `json:"bssid" yaml:"bssid" mapstructure:"bssid"`
}

type Configuration struct {
	Listener     ListenerEntry            `json:"listener" yaml:"listener" mapstructure:"listener"`
	Lockout      LockoutEntry             `json:"lockout" yaml:"lockout" mapstructure:"lockout"`
	Subsystems   SubsystemEntry           `json:"subsystems" yaml:"subsystems" mapstructure:"subsystems"`
	Updater      UpdaterEntry             `json:"updater" yaml:"updater" mapstructure:"updater"`
	WifiWatchdog WatchdogEntry            `json:"wifi_watchdog" yaml:"wifi_watchdog" mapstructure:"wifi_watchdog"`
	Plugins      PluginEntry              `json:"plugins" yaml:"plugins" mapstructure:"plugins"`
	Custom       []map[string]*RouteEntry `json:"routes" yaml:"routes" mapstructure:"routes"`
}

func Load(preferredLocation string) (err error) {
	// Setup configuration file location(s)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if preferredLocation != "" {
		viper.AddConfigPath(preferredLocation)
	}
	viper.AddConfigPath(GLOBAL_CONFIG_LOCATION)
	viper.AddConfigPath(LOCAL_CONFIG_LOCATION)
	viper.AddConfigPath(".")

	// Get the hostname and domain
	hostname, _ := os.Hostname()

	// Setup default values
	kvp := map[string]interface{}{
		"listener.port":                        8180,
		"listener.https.enabled":               true,
		"listener.https.type":                  "self-signed",
		"listener.https.create_if_missing":     true,
		"listener.https.domains":               []string{"localhost", "127.0.0.1", hostname},
		"listener.https.cert":                  DEFAULT_TLS_CERT_NAME,
		"listener.https.key":                   DEFAULT_TLS_KEY_NAME,
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
		"plugins.dir":                          DEFAULT_PLUGIN_LOCATION,
		"plugins.entries.hello.enabled":        true,
		"plugins.entries.hello.cookie:":        "this_is_not_a_security_feature",
		"plugins.entries.hello.hash":           "abc123",
		"plugins.entries.hello.config.message": "hello world plugin",
		"subsystems.iperf":                     true,
		"subsystems.tcp_reflector":             true,
		"subsystems.udp_reflector":             true,
		"subsystems.firewall":                  true,
		"routes":                               []map[string]*RouteEntry{},
	}
	for k, v := range kvp {
		viper.SetDefault(k, v)
	}

	if err := viper.ReadInConfig(); err != nil {
		// If the error is something other than the configuration file isn't found then
		// throw a fatal error for the user to handle.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read configuration file: %v", err)
		}

		// If the file is simply not found then create a default one and display a warning
		// to the user
		log.Println("[WARN] Configuration file not found")
		if ensureDir(GLOBAL_CONFIG_LOCATION, true) && checkWriteAccess(GLOBAL_CONFIG_LOCATION) {
			log.Printf("[INFO] Creating configuration file %sconfig.yaml", GLOBAL_CONFIG_LOCATION)
			if err := os.MkdirAll(GLOBAL_CONFIG_LOCATION, 0700); err != nil {
				log.Printf("[WARN] Error creating global config directory: %v\n", err)
			}
			err := viper.WriteConfigAs(path.Join(GLOBAL_CONFIG_LOCATION, "config.yaml"))
			if err != nil {
				return fmt.Errorf("failed to write log file: %v", err)
			}
		} else {
			location := strings.Replace(LOCAL_CONFIG_LOCATION, "$HOME", os.Getenv("HOME"), 1)
			if err := os.MkdirAll(location, 0700); err != nil {
				log.Printf("[WARN] Error creating user config directory: %v\n", err)
			}
			log.Printf("[INFO] Creating configuration file %s.config.yaml", location)
			err := viper.WriteConfigAs(path.Join(location, "config.yaml"))
			if err != nil {
				return fmt.Errorf("failed to write log file: %v", err)
			}
		}
		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("failed to read configuration file: %v", err)
		}
	}

	var c Configuration
	if err := viper.Unmarshal(&c); err != nil {
		return fmt.Errorf("Failed to unmarshal configuration: %v", err)
	}

	Config = &c

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
	return nil
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
	tempFile, err := ioutil.TempFile(dir, "write-check-")
	if err != nil {
		fmt.Printf("Error checking access: %v", err)
		return false
	}

	// If successful, delete the temporary file and return true.
	defer os.Remove(tempFile.Name())
	return true
}
