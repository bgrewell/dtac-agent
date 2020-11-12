package handlers

import (
	"fmt"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/network"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetInterfacesHandler(c *gin.Context) {
	ifaces, err := network.GetInterfaces()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, ifaces)
}

func GetInterfaceNamesHandler(c *gin.Context) {
	names, err := network.GetInterfaceNames()
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	WriteResponseJSON(c, names)
}

func GetInterfaceByNameHandler(c *gin.Context) {
	name := c.Param("name")
	if name != "" {
		var iface *network.Interface
		_, err := strconv.ParseInt(name, 10, 64)
		if err == nil {
			iface, err = network.GetInterfaceByIdx(name)
		} else {
			iface, err = network.GetInterfaceByName(name)
		}
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, iface)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving name"))
	}
}

func GetInterfaceByIdxHandler(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		iface, err := network.GetInterfaceByIdx(id)
		if err != nil {
			WriteErrorResponseJSON(c, err)
			return
		}
		WriteResponseJSON(c, iface)
	} else {
		WriteErrorResponseJSON(c, fmt.Errorf("error retrieving id"))
	}
}
