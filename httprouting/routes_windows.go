//go:build windows
// +build windows

package httprouting

import (
	"github.com/BGrewell/dtac-agent/handlers"
	"github.com/gin-gonic/gin"
)

func AddOSSpecificHandlers(r *gin.Engine) {
	// Wifi Watchdog
	//r.GET("/watchdog/wifi", handlers.GetWifiWatchdogHandler) // postman //TODO: Uncomment once watchdog fixed

	// NetQosPolicy Routes
	r.GET("/network/qos/policies", handlers.GetNetQosPoliciesHandler)         //Return all NetQosPolicy objects
	r.GET("/network/qos/policy/:name", handlers.GetNetQosPolicyHandler)       //Return named NetQosPolicy object
	r.PUT("/network/qos/policy/:name", handlers.UpdateNetQosPolicyHandler)    //Update a NetQosPolicy object
	r.POST("/network/qos/policy", handlers.CreateNetQosPolicyHandler)         //Create a new NetQosPolicy object
	r.DELETE("/network/qos/policy/:name", handlers.DeleteNetQosPolicyHandler) //Remove a NetQosPolicy object
	r.DELETE("/network/qos/policies", handlers.DeleteNetQosPoliciesHandler)   //Remove all NetQosPolicy objects
}
