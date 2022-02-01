package plugin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"plugin"
)

type Plugin interface {
	Register(config map[string]interface{}, r *gin.RouterGroup) error
	Name() string
}

func Initialize(dir string, config []map[string]map[string]interface{}, r *gin.Engine) {
	plugins, err := FindPlugins(dir, "*.so")
	if err != nil {
		log.Errorf("failed to find plugins directory: %s\n", err)
	}
	for _, pluginName := range plugins {

		library, err := plugin.Open(pluginName)
		if err != nil {
			log.Errorf("failed to load plugin: %s\n", err)
			continue
		}

		loader, err := library.Lookup("Load")
		if err != nil {
			log.Errorf("failed to lookup Load() function: %s\n", err)
			continue
		}

		plug := loader.(func() Plugin)()
		name := plug.Name()
		log.Infof("loading plugin %s\n", name)

		// Try to find a configuration for this plugin
		var c map[string]interface{}
		for _, pluginEntry := range config {
			for pluginName, pluginConfig := range pluginEntry {
				if pluginName == name {
					c = pluginConfig
					break
				}
			}
		}

		group := r.Group(fmt.Sprintf("plugins/%s", name))

		err = plug.Register(c, group)
		if err != nil {
			log.Errorf("failed to register plugin %s: %s\n", name, err)
			continue
		}
		log.Infof("registered plugin %s\n", name)
	}
}

func FindPlugins(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
