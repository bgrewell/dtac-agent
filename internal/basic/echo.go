package basic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"go.uber.org/zap"
)

// EchoArgs is a struct to assist with validating the input arguments
type EchoArgs struct {
	Message string `json:"msg"`
}

// EchoOutput is a struct to assist with validating the output
type EchoOutput struct {
	Message string `json:"message"`
}

// NewEchoSubsystem creates a new echo subsystem
func NewEchoSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "echo"
	es := EchoSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		enabled:    c.Config.Subsystems.Echo,
		name:       name,
	}
	es.register()
	return &es
}

// EchoSubsystem is a simple example subsystem for showing how the pieces fit together
type EchoSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	enabled    bool        // Optional subsystems have a boolean to control if they are enabled
	name       string      // Subsystem name
	endpoints  []*endpoint.Endpoint
}

// register registers the routes that this module handles
func (es *EchoSubsystem) register() {
	if !es.Enabled() {
		es.Logger.Info("subsystem is disabled", zap.String("subsystem", es.Name()))
		return
	}

	// Create a group for this subsystem
	base := es.name

	// Routes
	secure := es.Controller.Config.Auth.DefaultSecure
	authz := endpoint.AuthGroupAdmin.String()
	es.endpoints = []*endpoint.Endpoint{
		endpoint.NewEndpoint(fmt.Sprintf("%s/", base), endpoint.ActionRead, es.rootHandler, secure, authz, endpoint.WithParameters(EchoArgs{}), endpoint.WithOutput(EchoOutput{})),
	}
}

// Enabled returns whether the echo subsystem is enabled
func (es *EchoSubsystem) Enabled() bool {
	return es.enabled
}

// Name returns the name of the echo subsystem
func (es *EchoSubsystem) Name() string {
	return es.name
}

// Endpoints returns an array of endpoints that this Subsystem handles
func (es *EchoSubsystem) Endpoints() []*endpoint.Endpoint {
	return es.endpoints
}

func (es *EchoSubsystem) rootHandler(in *endpoint.EndpointRequest) (out *endpoint.EndpointResponse, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		if m := in.Parameters["msg"]; m[0] != "" {
			msg, err := json.Marshal(EchoOutput{Message: m[0]})
			return msg, err
		}

		return nil, errors.New("missing parameter 'msg'")

	}, "diagnostic information")
}
