package plugin

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	hclog "github.com/hashicorp/go-hclog"
	gplug "github.com/hashicorp/go-plugin"
)

var (
	binaryName string
)

func init() {
	if runtime.GOOS == "windows" {
		binaryName = "remote.exe"
	} else {
		binaryName = "remote"
	}
}

type RemoteConfiguration struct {
	HostAddr string
	HostPort int
}

// Controller hosts plugins and handles registering them with the System-API REST server
type Controller struct {
	HostAddr string
	HostPort int
	Plugins map[string]gplug.Plugin
	Clients map[string]*gplug.Client
}

func (c *Controller) DiscoverPlugins() error {
	// Discover the plugins available on the remote system
	args := []string{c.HostAddr, strconv.Itoa(c.HostPort), "list"}
	stdout, err := exec.Command(binaryName, args...).Output()
	if err != nil {
		return err
	}
	fmt.Println(stdout)
	c.Plugins = make(map[string]gplug.Plugin)
	pluginNames := strings.Split(string(stdout), "\n")
	for _, pluginName := range pluginNames {
		c.Plugins[pluginName] = &RestPlugin{}
	}
	return nil
}

func (c *Controller) LoadPlugins() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name: "plugin",
		Output: os.Stdout,
		Level: hclog.Debug,
	})

	for _, plugin := range c.Plugins {
		client := gplug.NewClient(&gplug.ClientConfig{
			HandshakeConfig: gplug.HandshakeConfig{},
			Plugins:         nil,
			Cmd:             nil,
			Logger:          nil,
		})
	}
}
