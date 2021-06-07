package handlers

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/configuration"
	"github.com/BGrewell/wifi-watchdog"
	"time"
)

var (
	Watchdog wifi-watchdog.WifiWatchdog
)

func init() {
	go func() {
		// wait for the config to be loaded then create the watchdog and start it
		for {
			cfg, err := configuration.GetActiveConfig()
			if err == nil {

			}
		}
	}()
}

func GetWifiWatchdogHandler(c *gin.Context) {
	start := time.Now()
	WriteResponseJSON(c, time.Since(start), "this method has not been implemented. the watchdog is enabled but display/modifications are not.")
}
