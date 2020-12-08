// +build windows

package httprouting

import (
	"github.com/BGrewell/system-api/handlers"
	"github.com/gin-gonic/gin"
)

func AddOSSpecificHandlers(r *gin.Engine) {
	// NetQosPolicy Routes
	r.GET("/network/qos/policies", handlers.GetNetQosPolicies)         //Return all NetQosPolicy objects
	r.GET("/network/qos/policy/:name", handlers.GetNetQosPolicy)       //Return named NetQosPolicy object
	r.PUT("/network/qos/policy/:name", handlers.UpdateNetQosPolicy)    //Update a NetQosPolicy object
	r.POST("/network/qos/policy", handlers.CreateNetQosPolicy)         //Create a new NetQosPolicy object
	r.DELETE("/network/qos/policy/:name", handlers.DeleteNetQosPolicy) //Remove a NetQosPolicy object
	r.DELETE("/network/qos/policies", handlers.DeleteNetQosPolicies)   //Remove all NetQosPolicy objects
}
