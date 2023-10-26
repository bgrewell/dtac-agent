package helpers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
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
func (f *JSONResponseFormatter) WriteResponse(c *gin.Context, duration time.Duration, obj interface{}) {
	var response interface{}

	rawData, err := json.Marshal(obj)
	if err != nil {
		f.WriteError(c, err)
		return
	}
	rawMsg := json.RawMessage(rawData)

	// We ensure that the response is always wrapped in a ResponseWrapper struct in order to ensure that the response
	// is always properly formatted as machine-readable input, for example if a plugin or other function returns a
	// string.
	// WrapResponses is a configuration option that wraps the response with additional information including timing
	// information and a text based status code. This is useful for debugging and testing purposes.
	if f.cfg.Output.WrapResponses {
		response = OKResponse{
			Time:    time.Now().Format(time.RFC3339Nano),
			Status:  "success",
			Elapsed: duration.String(),
			ResponseWrapper: ResponseWrapper{
				Response: rawMsg,
			},
		}
	} else {
		response = ResponseWrapper{
			Response: rawMsg,
		}
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
	c.Data(http.StatusInternalServerError, gin.MIMEJSON, jerr)
	c.Abort()
}

// WriteNotImplementedError writes a not implemented error response in JSON format
func (f *JSONResponseFormatter) WriteNotImplementedError(c *gin.Context, err error) {
	er := ErrorResponse{
		Time: time.Now().Format(time.RFC3339Nano),
		Err:  err.Error(),
	}
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
	jerr, _ := json.Marshal(er)
	c.Data(http.StatusNotFound, gin.MIMEJSON, jerr)
	c.Abort()
}
