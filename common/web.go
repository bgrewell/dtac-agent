package common

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func WriteResponseJSON(c *gin.Context, obj interface{}) {
	jout, err := json.Marshal(obj)
	if err != nil {
		WriteErrorResponseJSON(c, err)
	}
	c.Data(http.StatusOK, gin.MIMEJSON, jout)
}

type ErrorResponse struct {
	Err  string `json:"error"`
	Time string `json:"time"`
}

func WriteErrorResponseJSON(c *gin.Context, err error) {
	er := ErrorResponse{
		Err:  err.Error(),
		Time: time.Now().Format(time.RFC3339Nano),
	}
	jerr, err := json.Marshal(er)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
	}
	c.Data(http.StatusInternalServerError, gin.MIMEJSON, jerr)
}
