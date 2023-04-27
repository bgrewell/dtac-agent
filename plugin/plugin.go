package plugin

import (
	"github.com/bgrewell/gin-plugins/loader"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Initialize(dir string, cookie string, config []map[string]map[string]interface{}, r *gin.Engine) (err error) {
	group := r.Group("plugins")

	l := loader.NewPluginLoader(dir, cookie, group)
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
