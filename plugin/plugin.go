package plugin

import (
	"fmt"
	"github.com/BGrewell/system-api/configuration"
	"github.com/gin-gonic/gin"
)

var (
	RegisteredPlugins map[string]Plugin
)

func init() {
	RegisteredPlugins = make(map[string]Plugin)
}

type Plugin interface {
	Register(r *gin.Engine)
	Name() string
}

func Initialize(config configuration.Plugins) {
	for idx, plugin := range config.Entries {
		fmt.Println("initializing %s: %v\n", idx, plugin)
		if plug, ok := RegisteredPlugins[plugin]; ok {

		}
	}
}
