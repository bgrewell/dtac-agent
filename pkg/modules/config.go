package modules

import "path/filepath"

// ModuleConfig is the configuration for a module
type ModuleConfig struct {
	ModulePath string                 `json:"module_path" yaml:"module_path"`
	RootPath   string                 `json:"root_path,omitempty" yaml:"root_path,omitempty"`
	Enabled    bool                   `json:"enabled" yaml:"enabled"`
	Hash       string                 `json:"hash" yaml:"hash"`
	User       string                 `json:"user" yaml:"user"`
	Config     map[string]interface{} `json:"config" yaml:"config"`
}

// Name returns the name of the module
func (mc ModuleConfig) Name() string {
	return filepath.Base(mc.ModulePath)
}
