package rest

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

// OKResponse is the struct for the response
type OKResponse struct {
	Time    string `json:"time" yaml:"time" xml:"time"`
	Status  string `json:"status" yaml:"status" xml:"status"`
	Elapsed string `json:"processing_time" yaml:"elapsed" xml:"elapsed"`
	ResponseWrapper
}

// ResponseWrapper is a struct to ensure that all output always is properly formatted as machine-readable input
type ResponseWrapper struct {
	Response json.RawMessage `json:"response" yaml:"response" xml:"response"`
}

// ErrorResponse is the struct for the response
type ErrorResponse struct {
	Time string `json:"time"`
	Err  string `json:"error"`
}

// ResponseFormatter is the interface for the response formatters used to write responses to the API service(s)
type ResponseFormatter interface {
	WriteResponse(c *gin.Context, duration time.Duration, obj []byte)
	WriteError(c *gin.Context, err error)
	WriteNotImplementedError(c *gin.Context, err error)
	WriteUnauthorizedError(c *gin.Context, err error)
	WriteNotFoundError(c *gin.Context)
}
