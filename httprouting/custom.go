package httprouting

import (
	"fmt"
	"github.com/BGrewell/dtac-agent/configuration"
	"github.com/BGrewell/dtac-agent/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AddCustomHandlers(c *configuration.Config, r *gin.Engine) {

	for _, entry := range c.Custom {
		for _, value := range entry {
			route := fmt.Sprintf("/custom/%s", value.Name)
			if value.Route != "" {
				route = value.Route
			}
			//TODO: Check for existing key in map, if found warn user that the existing route is being shadowed
			settings := &handlers.CustomHandlerSettings{
				Value:    value.Source.Value,
				Type:     value.Source.Type,
				Settings: value.Blocking,
			}
			handlers.AddCustomHandler(route, settings)

			switch value.Source.Type {
			case "file":
				r.GET(route, handlers.CustomFileHandler)
			case "net":
				r.GET(route, handlers.CustomNetHandler)
			case "cmd":
				r.GET(route, handlers.CustomCmdHandler)
			default:
				// todo: need to push logging throughout the code
				log.Printf("unrecognized custom route source: %v. skipping", value.Source.Type)
			}
		}
	}
}
