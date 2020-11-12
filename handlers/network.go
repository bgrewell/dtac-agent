package handlers

import (
	"fmt"
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/network"
	"github.com/gorilla/mux"
	"net/http"
)

func GetInterfacesHandler(w http.ResponseWriter, r *http.Request) {
	ifaces, err := network.GetInterfaces()
	if err != nil {
		WriteErrorResponseJSON(w, err)
		return
	}
	WriteResponseJSON(w, ifaces)
}

func GetInterfaceNamesHandler(w http.ResponseWriter, r *http.Request) {
	names, err := network.GetInterfaceNames()
	if err != nil {
		WriteErrorResponseJSON(w, err)
		return
	}
	WriteResponseJSON(w, names)
}

func GetInterfaceByNameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if val, ok := vars["name"]; ok {
		iface, err := network.GetInterfaceByName(val)
		if err != nil {
			WriteErrorResponseJSON(w, err)
			return
		}
		WriteResponseJSON(w, iface)
	} else {
		WriteErrorResponseJSON(w, fmt.Errorf("error retrieving name: %v", vars))
	}
}

func GetInterfaceByIdxHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if val, ok := vars["id"]; ok {
		iface, err := network.GetInterfaceByIdx(val)
		if err != nil {
			WriteErrorResponseJSON(w, err)
			return
		}
		WriteResponseJSON(w, iface)
	} else {
		WriteErrorResponseJSON(w, fmt.Errorf("error retrieving id: %v", vars))
	}
}