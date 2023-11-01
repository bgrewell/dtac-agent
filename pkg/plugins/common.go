package plugins

import "encoding/json"

// PluginCommon is a common struct that holds the common methods. It must be seperate from
// the PluginBase struct in order to prevent issues with the net.rpc trying to
// serve functions that don't fit the rpc function signatures.
type PluginCommon struct {
}

// Name returns the name of the plugin
func (p *PluginCommon) Name() string {
	return "UnnamedPlugin"
}

// RootPath returns the root path for the plugin
func (p *PluginCommon) RootPath() string {
	return ""
}

// Serialize serializes the given interface to a string
func (p *PluginCommon) Serialize(v interface{}) (string, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return "", e
	}
	return string(b), nil
}
