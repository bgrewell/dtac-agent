package plugins

import "path/filepath"

// PluginConfig is the configuration for a plugin
type PluginConfig struct {
	PluginPath string                 `json:"plugin_path" yaml:"plugin_path"`
	RootPath   string                 `json:"root_path,omitempty" yaml:"root_path,omitempty"`
	Enabled    bool                   `json:"enabled" yaml:"enabled"`
	Hash       string                 `json:"hash" yaml:"hash"`
	User       string                 `json:"user" yaml:"user"`
	Config     map[string]interface{} `json:"config" yaml:"config"`
}

// Name returns the name of the plugin
func (pc PluginConfig) Name() string {
	return filepath.Base(pc.PluginPath)
}
