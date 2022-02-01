package plugin

import (
	"github.com/BGrewell/system-api/plugin/hello"
	"github.com/BGrewell/system-api/plugin/remote"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	RegisteredPlugins map[string]Plugin
)

func init() {
	RegisteredPlugins = map[string]Plugin{
		"hello":  &hello.HelloPlugin{},
		"remote": &remote.RemotePlugin{},
	}
}

type Plugin interface {
	Register(config map[string]interface{}, r *gin.RouterGroup) error
	Name() string
}

func Initialize(config []map[string]map[string]interface{}, r *gin.Engine) {
	for _, plugEntry := range config {
		for pluginName, pluginConfig := range plugEntry {
			if plugin, ok := RegisteredPlugins[pluginName]; ok {
				plugRouter := r.Group(plugin.Name())
				err := plugin.Register(pluginConfig, plugRouter)
				if err != nil {
					log.Errorf("failed to register plugin %s: %v\n", plugin.Name(), err)
				}
			} else {
				log.Errorf("Unknown plugin: %s. Configuration ignored!", pluginName)
			}
		}
	}
}
