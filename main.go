package main

import (
	. "github.com/BGrewell/system-api/common"
	"github.com/BGrewell/system-api/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type HomeResponse struct {
	Status string `json:"system_status"`
	Time string `json:"request_time"`
	Routes []string `json:"routes"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	h := &HomeResponse{
		Status: "OK",
		Time:   time.Now().Format(time.RFC3339Nano),
		Routes: []string{
			"/",
			"/network/interfaces",
			"/network/interfaces/names",
			"/network/interface/{name:str}",
			"/network/interface/idx/{idx:int}",
		},
	}
	WriteResponseJSON(w, h)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/network/interfaces", handlers.GetInterfacesHandler)
	router.HandleFunc("/network/interfaces/names", handlers.GetInterfaceNamesHandler)
	router.HandleFunc("/network/interface/{name}", handlers.GetInterfaceByNameHandler)
	router.HandleFunc("/network/interface/idx/{idx:[0-9]+}", handlers.GetInterfaceByIdxHandler)

	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
