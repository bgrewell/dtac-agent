package plugin

import (
	"github.com/BGrewell/system-api/configuration"
	"strings"
)

type Client struct {
	PluginName string
	launcher   Launcher
	cfg        *configuration.PluginEntry
}

/*
NewClient creates a new plugin client

Creating a new client includes deploying and launching the plugin and creating the grpc client connection.
*/
func NewClient(pluginName string, cfg *configuration.PluginEntry) (client *Client, err error) {

	var l Launcher
	switch strings.ToLower(cfg.Protocol) {
	case "ssh":
		l = SSHLauncher{}
	case "local":
		l = LocalLauncher{}
	}

	c := &Client{
		PluginName: pluginName,
		launcher:   l,
		cfg:        cfg,
	}

	return c, nil
}
