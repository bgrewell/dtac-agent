package main

import (
	"encoding/json"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/network"
	"time"
)

func main() {
	stats, _ := network.GetInterfaceStats("LAN2")
	time.Sleep(1 * time.Second)
	stats, err := network.GetInterfaceStats("LAN2")
	if err != nil {
		fmt.Println(err)
		return
	}
	statsout, err := json.Marshal(stats)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(statsout))
}
