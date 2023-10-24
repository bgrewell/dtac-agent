package helpers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	MOVE_ME_TO_CONFIG_SHOW_WRAPPED_OUTPUT = false
)

// OKResponse is the struct for the response
type OKResponse struct {
	Time    string      `json:"time"`
	Status  string      `json:"status"`
	Elapsed string      `json:"processing_time"`
	Output  interface{} `json:"output"`
}

// ErrorResponse is the struct for the response
type ErrorResponse struct {
	Time string `json:"time"`
	Err  string `json:"error"`
}

// WriteResponseJSON writes a response in JSON format
func WriteResponseJSON(c *gin.Context, duration time.Duration, obj interface{}) {
	response := obj
	if MOVE_ME_TO_CONFIG_SHOW_WRAPPED_OUTPUT {
		response = OKResponse{
			Time:    time.Now().Format(time.RFC3339Nano),
			Status:  "success",
			Elapsed: duration.String(),
			Output:  obj,
		}
	}
	jout, err := json.Marshal(response)
	if err != nil {
		WriteErrorResponseJSON(c, err)
		return
	}
	c.Data(http.StatusOK, gin.MIMEJSON, jout)
}

// WriteErrorResponseJSON writes an error response in JSON format
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
	c.Abort()
}

// WriteUnauthorizedResponseJSON writes an error response in JSON format
func WriteUnauthorizedResponseJSON(c *gin.Context, err error) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  err.Error(),
	}
	jerr, err := json.Marshal(er)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}
	c.Data(http.StatusUnauthorized, gin.MIMEJSON, jerr)
	c.Abort()
}

// WriteNotFoundResponseJSON writes an error response in JSON format
func WriteNotFoundResponseJSON(c *gin.Context) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  "404 page not found",
	}
	jerr, _ := json.Marshal(er)
	c.Data(http.StatusNotFound, gin.MIMEJSON, jerr)
	c.Abort()
}

// WriteNotImplementedResponseJSON writes an error response in JSON format
func WriteNotImplementedResponseJSON(c *gin.Context) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  "this method has not been implemented yet",
	}
	jerr, err := json.Marshal(er)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}
	c.Data(http.StatusNotImplemented, gin.MIMEJSON, jerr)
	c.Abort()
}
