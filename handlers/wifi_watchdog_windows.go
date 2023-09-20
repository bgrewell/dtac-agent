package handlers

// TODO: This is commented out because of import issues with BGrewell/wifi-watchdog being private.

//import (
//	"github.com/intel-innersource/frameworks.automation.dtac.agent/configuration"
//	"github.com/gin-gonic/gin"
//	log "github.com/sirupsen/logrus"
//	"time"
//)
//
//var (
//	Watchdog watchdog.WifiWatchdog
//)
//
////TODO: The watchdog should be a module or even a plugin not here. Refactor this before releasing
//
//func init() {
//	go func() {
//		// wait for the config to be loaded then create the watchdog and start it
//		for {
//			cfg, err := configuration.GetActiveConfig()
//			if err == nil {
//				Watchdog = watchdog.WifiWatchdog{}
//				if cfg.Watchdog.Enabled {
//					Watchdog.Initialize()
//					Watchdog.Start(cfg.Watchdog.Profile, cfg.Watchdog.PollInterval)
//					log.Printf("wifi watchdog started. watching for profile: %s every %d seconds\n", cfg.Watchdog.Profile, cfg.Watchdog.PollInterval)
//				}
//				return
//			}
//			time.Sleep(30 * time.Second)
//		}
//
//	}()
//}
//
//func GetWifiWatchdogHandler(c *gin.Context) {
//	start := time.Now()
//	WriteResponseJSON(c, time.Since(start), "this method has not been implemented. the watchdog is enabled but display/modifications are not.")
//}
