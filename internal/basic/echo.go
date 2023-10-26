package basic

import (
	"errors"
	"fmt"

	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"go.uber.org/zap"
)

// EchoArgs is a struct to assist with validating the input arguments
type EchoArgs struct {
	Message string `json:"msg"`
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
	endpoints  []endpoint.Endpoint
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
	es.endpoints = []endpoint.Endpoint{
		{Path: fmt.Sprintf("%s/", base), Action: endpoint.ActionRead, Function: es.rootHandler, UsesAuth: secure, ExpectedArgs: EchoArgs{}, ExpectedBody: nil},
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
func (es *EchoSubsystem) Endpoints() []endpoint.Endpoint {
	return es.endpoints
}

func (es *EchoSubsystem) rootHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		if m := in.Params["msg"]; m[0] != "" {
			msg := m[0]
			return msg, nil
		}

		return nil, errors.New("missing parameter 'msg'")

	}, "diagnostic information")
}
