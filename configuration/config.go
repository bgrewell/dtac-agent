package configuration

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type WatchdogEntry struct {
	Enabled      bool   `json:"enabled" yaml:"enabled" xml:"enabled"`
	PollInterval int    `json:"poll_interval" yaml:"poll_interval" xml:"poll_interval"`
	Profile      string `json:"profile" yaml:"profile" xml:"profile"`
	BSSID        string `json:"bssid" yaml:"bssid" xml:"bssid"`
}

type SourceEntry struct {
	Type  string `json:"type" yaml:"type" xml:"type"`
	Value string `json:"value" yaml:"value" xml:"value"`
	RunAs string `json:"run_as" yaml:"run_as" xml:"run_as"`
}

type BlockingEntry struct {
	Trigger       string `json:"trigger" yaml:"trigger" xml:"trigger"`
	Detect        string `json:"detect" yaml:"detect" xml:"detect"`
	TimeoutMs     int    `json:"timeout_ms" yaml:"timeout_ms" xml:"timeout_ms"`
	TimeoutAction string `json:"timeout_action" yaml:"timeout_action" xml:"timeout_action"`
}

type CustomEntry struct {
	Name     string         `json:"name" yaml:"name" xml:"name"`
	Route    string         `json:"route" yaml:"route" xml:"route"`
	Source   *SourceEntry   `json:"source" yaml:"source" xml:"source"`
	Methods  []string       `json:"methods" yaml:"methods" xml:"methods"`
	Access   string         `json:"access" yaml:"access" xml:"access"`
	Blocking *BlockingEntry `json:"blocking" yaml:"blocking" xml:"blocking"`
	WrapJson bool           `json:"wrap_json" yaml:"wrap_json" xml:"wrap_json"`
	Mode     string         `json:"mode" yaml:"mode" xml:"mode"`
}

type UpdaterEntry struct {
	Token           string `json:"token" yaml:"token" xml:"token"`
	Mode            string `json:"mode" yaml:"mode" xml:"mode"`
	Interval        string `json:"interval" yaml:"interval" xml:"interval"`
	ErrorFallback   string `json:"error_fallback" yaml:"error_fallback" xml:"error_fallback"`
	RestartOnUpdate bool   `json:"restart_on_update" yaml:"restart_on_update" xml:"restart_on_update"`
}

//type PluginsEntry struct {
//	ListenPort    int                       `json:"listen_port" yaml:"listen_port" xml:"listen_port"`
//	PluginDir     string                    `json:"dir" yaml:"dir" xml:"dir"`
//	ActivePlugins []map[string]*PluginEntry `json:"active" yaml:"active" xml:"active"`
//}

//type PluginEntry struct {
//	Binary        string `json:"binary" yaml:"binary" xml:"binary"`
//	Target        string `json:"target" yaml:"target" xml:"target"`
//	Protocol      string `json:"protocol" yaml:"protocol" xml:"protocol"`
//	User          string `json:"user" yaml:"user" xml:"user"`
//	Pass          string `json:"pass" yaml:"pass" xml:"pass"`
//	EnsureRunning bool   `json:"ensure_running" yaml:"ensure_running" xml:"ensure_running"`
//}

type Config struct {
	ListenPort  int                                 `json:"listen_port" yaml:"listen_port" xml:"listen_port"`
	HTTPS       bool                                `json:"https" yaml:"https" xml:"https"`
	CertFile    string                              `json:"cert_file" yaml:"cert_file" xml:"cert_file"`
	KeyFile     string                              `json:"key_file" yaml:"key_file" xml:"key_file"`
	LockoutTime int                                 `json:"lockout_timeout" yaml:"lockout_timeout" xml:"lockout_time"`
	Updater     UpdaterEntry                        `json:"updater" yaml:"updater" xml:"updater"`
	PluginDir   string                              `json:"plugin_dir" yaml:"plugin_dir" xml:"plugin_dir"`
	Plugins     []map[string]map[string]interface{} `json:"plugins" yaml:"plugins" xml:"plugins"`
	Modules     []map[string]map[string]interface{} `json:"modules" yaml:"modules" xml:"modules"`
	Custom      []map[string]*CustomEntry           `json:"custom" yaml:"custom" xml:"custom"`
	Watchdog    WatchdogEntry                       `json:"watchdog" yaml:"watchdog" xml:"watchdog"`
}

var (
	instance *Config
)

func GetActiveConfig() (config *Config, err error) {
	if instance == nil {
		return nil, errors.New("no active configuration exists")
	}

	return instance, nil
}

func Load(filename string) (config *Config, err error) {
	c := &Config{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	// populate names
	for _, entry := range c.Custom {
		for key, value := range entry {
			value.Name = key
		}
	}
	if c.LockoutTime == 0 {
		c.LockoutTime = 60
	}
	instance = c
	return c, nil
}
