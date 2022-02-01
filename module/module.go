package module

import (
	"github.com/BGrewell/system-api/module/hello"
	"github.com/BGrewell/system-api/module/remote"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	RegisteredModules map[string]Module
)

func init() {
	RegisteredModules = map[string]Module{
		"hello":  &hello.HelloModule{},
		"remote": &remote.RemoteModule{},
	}
}

type Module interface {
	Register(config map[string]interface{}, r *gin.RouterGroup) error
	Name() string
}

func Initialize(config []map[string]map[string]interface{}, r *gin.Engine) {
	for _, moduleEntry := range config {
		for moduleName, moduleConfig := range moduleEntry {
			if module, ok := RegisteredModules[moduleName]; ok {
				plugRouter := r.Group(module.Name())
				err := module.Register(moduleConfig, plugRouter)
				if err != nil {
					log.Errorf("failed to register module %s: %v\n", module.Name(), err)
				}
			} else {
				log.Errorf("Unknown module: %s. Configuration ignored!", moduleName)
			}
		}
	}
}
