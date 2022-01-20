package core

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type OKResponse struct {
	Time    string      `json:"time"'`
	Status  string      `json:"status"`
	Elapsed string      `json:"processing_time"`
	Output  interface{} `json:"output"`
}

func WriteResponseJSON(c *gin.Context, duration time.Duration, obj interface{}) {
	response := OKResponse{
		Time:    time.Now().Format(time.RFC3339Nano),
		Status:  "success",
		Elapsed: fmt.Sprintf("%s", duration),
		Output:  obj,
	}
	jout, err := json.Marshal(response)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	c.Data(http.StatusOK, gin.MIMEJSON, jout)
}

type ErrorResponse struct {
	Time string `json:"time"`
	Err  string `json:"error"`
}

func WriteErrorResponseJSON(c *gin.Context, err error) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  err.Error(),
	}
	jerr, err := json.Marshal(er)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}
	c.Data(http.StatusInternalServerError, gin.MIMEJSON, jerr)
}
