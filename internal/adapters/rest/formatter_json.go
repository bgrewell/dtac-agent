package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/bgrewell/dtac-agent/internal/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// NewJSONResponseFormatter creates a new instance of the JSON response formatter
func NewJSONResponseFormatter(cfg *config.Configuration, logger *zap.Logger) ResponseFormatter {
	name := "json_formatter"
	return &JSONResponseFormatter{
		cfg:    cfg,
		logger: logger.With(zap.String("module", name)),
	}
}

// JSONResponseFormatter is the struct for the JSON response formatter
type JSONResponseFormatter struct {
	cfg    *config.Configuration
	logger *zap.Logger
}

// WriteResponse writes a response in JSON format
func (f *JSONResponseFormatter) WriteResponse(c *gin.Context, duration time.Duration, obj []byte) {
	var response interface{}

	var js json.RawMessage
	if err := json.Unmarshal(obj, &js); err != nil {
		f.WriteError(c, err)
		return
	}

	c.Header("X-DTAC-Duration", duration.String())
	c.Header("X-DTAC-Status", "success")
	c.Header("X-DTAC-Time", time.Now().Format(time.RFC3339Nano))

	response = ResponseWrapper{
		Response: js,
	}

	jout, err := json.Marshal(response)
	if err != nil {
		f.WriteError(c, err)
		return
	}
	c.Data(http.StatusOK, gin.MIMEJSON, jout)
}

// WriteError writes an error response in JSON format
func (f *JSONResponseFormatter) WriteError(c *gin.Context, err error) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  err.Error(),
	}
	jerr, err := json.Marshal(er)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}

	c.Header("X-Exec-Status", "error")
	c.Header("X-Exec-Time", time.Now().Format(time.RFC3339Nano))

	c.Data(http.StatusInternalServerError, gin.MIMEJSON, jerr)
	c.Abort()
}

// WriteNotImplementedError writes a not implemented error response in JSON format
func (f *JSONResponseFormatter) WriteNotImplementedError(c *gin.Context, err error) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  err.Error(),
	}

	c.Header("X-Exec-Status", "not-implemented")
	c.Header("X-Exec-Time", time.Now().Format(time.RFC3339Nano))

	jerr, err := json.Marshal(er)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}
	c.Data(http.StatusNotImplemented, gin.MIMEJSON, jerr)
	c.Abort()
}

// WriteUnauthorizedError writes an unauthorized error response in JSON format
func (f *JSONResponseFormatter) WriteUnauthorizedError(c *gin.Context, err error) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  err.Error(),
	}

	c.Header("X-Exec-Status", "unauthorized")
	c.Header("X-Exec-Time", time.Now().Format(time.RFC3339Nano))

	jerr, err := json.Marshal(er)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "time": time.Now().Format(time.RFC3339Nano)})
		return
	}
	c.Data(http.StatusUnauthorized, gin.MIMEJSON, jerr)
	c.Abort()
}

// WriteNotFoundError writes a not found error response in JSON format
func (f *JSONResponseFormatter) WriteNotFoundError(c *gin.Context) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  "404 page not found",
	}

	c.Header("X-Exec-Status", "not-found")
	c.Header("X-Exec-Time", time.Now().Format(time.RFC3339Nano))

	jerr, _ := json.Marshal(er)
	c.Data(http.StatusNotFound, gin.MIMEJSON, jerr)
	c.Abort()
}
