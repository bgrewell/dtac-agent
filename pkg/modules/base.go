package modules

import (
	"errors"
	"fmt"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"log"
	"strings"
)

// LoggingLevel is an enum for the various logging levels
type LoggingLevel int

const (
	LoggingLevelDebug LoggingLevel = iota
	LoggingLevelInfo
	LoggingLevelWarning
	LoggingLevelError
	LoggingLevelFatal
)

// LogMessage represents a structured log message
type LogMessage struct {
	Level   LoggingLevel
	Message string
	Fields  map[string]string
}

// ModuleBase is a base struct that all modules should embed as it implements the common shared methods
type ModuleBase struct {
	LogChan        chan LogMessage
	Methods        map[string]endpoint.Func
	rootPath       string
	standaloneMode bool
}

// Register is a default implementation of the Register method that must be implemented by the module therefor this one returns an error
func (m *ModuleBase) Register(request *api.ModuleRegisterRequest, reply *api.ModuleRegisterResponse) error {
	return errors.New("this method must be implemented by the module")
}

// Call is a shim that calls the appropriate method on the module
func (m *ModuleBase) Call(method string, args *endpoint.Request) (out *endpoint.Response, err error) {
	key := method
	if f, exists := m.Methods[key]; exists {
		return f(args)
	}

	return nil, fmt.Errorf("method %s not found", method)
}

// LoggingStream is a function that sets up the logging channel for modules to use so that they can log messages back
// to the agent. In advanced cases this could be overridden by the module to implement its own handling of the logging
// stream but there likely isn't a good reason to do that.
func (m *ModuleBase) LoggingStream(stream api.ModuleService_LoggingStreamServer) error {
	if m.LogChan == nil {
		m.LogChan = make(chan LogMessage, 4096)
	}
	for {
		msg := <-m.LogChan
		fields := make([]*api.LogField, 0)
		for k, v := range msg.Fields {
			fields = append(fields, &api.LogField{
				Key:   k,
				Value: v,
			})
		}
		err := stream.Send(&api.LogMessage{
			Level:   api.LogLevel(msg.Level),
			Message: msg.Message,
			Fields:  fields,
		})
		if err != nil {
			return err
		}
	}
}

// RegisterMethods is used to create the call map for the module
func (m *ModuleBase) RegisterMethods(endpoints []*endpoint.Endpoint) {
	if m.Methods == nil {
		m.Methods = make(map[string]endpoint.Func)
	}
	for _, ep := range endpoints {
		m.Methods[fmt.Sprintf("%s:%s", ep.Action, ep.Path)] = ep.Function
	}
}

// Name returns the name of the module
func (m *ModuleBase) Name() string {
	return "UnnamedModule"
}

// RootPath returns the root path for the module
func (m *ModuleBase) RootPath() string {
	return m.rootPath
}

// SetRootPath sets the value of rootPath for the module
func (m *ModuleBase) SetRootPath(rootPath string) {
	m.rootPath = rootPath
}

// Log logs a message to the logging channel or stdout in standalone mode
func (m *ModuleBase) Log(level LoggingLevel, message string, fields map[string]string) {
	if m.standaloneMode {
		// Log to stdout in standalone mode
		levelStr := ""
		switch level {
		case LoggingLevelDebug:
			levelStr = "DEBUG"
		case LoggingLevelInfo:
			levelStr = "INFO"
		case LoggingLevelWarning:
			levelStr = "WARN"
		case LoggingLevelError:
			levelStr = "ERROR"
		case LoggingLevelFatal:
			levelStr = "FATAL"
		}
		
		if len(fields) > 0 {
			// Format fields as key=value pairs with pre-allocated capacity
			fieldStrs := make([]string, 0, len(fields))
			for k, v := range fields {
				// Escape newlines and tabs to prevent log injection
				k = strings.ReplaceAll(k, "\n", "\\n")
				k = strings.ReplaceAll(k, "\r", "\\r")
				k = strings.ReplaceAll(k, "\t", "\\t")
				v = strings.ReplaceAll(v, "\n", "\\n")
				v = strings.ReplaceAll(v, "\r", "\\r")
				v = strings.ReplaceAll(v, "\t", "\\t")
				fieldStrs = append(fieldStrs, fmt.Sprintf("%s=%s", k, v))
			}
			log.Printf("[%s] %s {%s}\n", levelStr, message, strings.Join(fieldStrs, ", "))
		} else {
			log.Printf("[%s] %s\n", levelStr, message)
		}
		return
	}
	
	// Use RPC logging when running under DTAC
	if m.LogChan == nil {
		m.LogChan = make(chan LogMessage, 4096)
	}
	m.LogChan <- LogMessage{
		Level:   level,
		Message: message,
		Fields:  fields,
	}
}

// SetStandaloneMode sets whether the module is running in standalone mode
func (m *ModuleBase) SetStandaloneMode(standalone bool) {
	m.standaloneMode = standalone
}
