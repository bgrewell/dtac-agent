package plugins

import "errors"

// PluginBase is a base struct that all plugins should embed as it implements the common shared methods
type PluginBase struct {
	PluginCommon
}

func (p *PluginBase) Register(args RegisterArgs, reply *RegisterReply) error {
	return errors.New("this method must be implemented by the plugin")
}
