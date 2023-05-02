package plugin

import (
	"github.com/bgrewell/gin-plugins/loader"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"path"
)

func Initialize(dir string, config map[string]*loader.PluginConfig, r *gin.Engine) (err error) {
	group := r.Group("plugins")

	// Remap the plugin configs to use full path for key
	cm := make(map[string]*loader.PluginConfig, 0)
	for k, v := range config {

		// Deal with any poorly formed entries
		if v == nil {
			log.Printf("bad plugin entry for %s\n", k)
			continue
		}
		full := path.Join(dir, k)
		v.PluginPath = full
		cm[full] = v
	}

	l := loader.NewPluginLoader(dir, cm, group, false)
	active, err := l.Initialize()
	if err != nil {
		return err
	}

	log.Printf("loaded %d plugins\n", len(active))
	for idx, plug := range active {
		log.Printf("  %d: %s [%s]\n", idx+1, plug.Name, plug.Path)
	}

	return nil
}
