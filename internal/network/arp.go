package network

import (
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"time"
)

// ArpEntry is the struct for the arp entry
type ArpEntry struct {
	IPAddress string `json:"ip_address"`
	HWType    string `json:"hw_type"`
	Flags     string `json:"flags"`
	HWAddress string `json:"hw_address"`
	Mask      string `json:"mask"`
	Iface     string `json:"device"`
}

func arpTableHandler(c *gin.Context) {
	start := time.Now()
	arpData, err := GetArpTable()
	if err != nil {
		helpers.WriteErrorResponseJSON(c, err)
		return
	}
	response := gin.H{
		"arp-table": types.AnnotatedStruct{
			Description: "returns ARP table information from the system",
			Value:       arpData,
		},
	}
	helpers.WriteResponseJSON(c, time.Since(start), response)
}
