package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func WriteResponseJSON(w http.ResponseWriter, obj interface{}){
	jout, err := json.Marshal(obj)
	if err != nil {
		WriteErrorResponseJSON(w, err)
	}
	_, _ = w.Write(jout)
}

type ErrorResponse struct {
	Err string `json:"error"`
	Time string `json:"time"`
}

func WriteErrorResponseJSON(w http.ResponseWriter, err error) {
	er := ErrorResponse{
		Err:  err.Error(),
		Time: time.Now().Format(time.RFC3339Nano),
	}
	w.WriteHeader(http.StatusInternalServerError)
	jerr, err := json.Marshal(er)
	if err != nil {
		_, _ = fmt.Fprintf(w, fmt.Sprintf("{error: %s}", err))
	}
	_, _ = w.Write(jerr)
}