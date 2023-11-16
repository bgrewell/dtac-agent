package basic

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/controller"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/version"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"go.uber.org/zap"
)

// HomeOutput is the struct for the home page output
type HomeOutput struct {
	Message   string               `json:"message"`
	Version   string               `json:"version"`
	Endpoints []*endpoint.Endpoint `json:"endpoints"`
}

// NewHomePageSubsystem creates a new homepage subsystem
func NewHomePageSubsystem(c *controller.Controller) interfaces.Subsystem {
	name := "homepage"
	hps := HomePageSubsystem{
		Controller: c,
		Logger:     c.Logger.With(zap.String("module", name)),
		name:       name,
	}
	hps.register()
	return &hps
}

// HomePageSubsystem serves the main homepage content
type HomePageSubsystem struct {
	Controller *controller.Controller
	Logger     *zap.Logger // All subsystems have a pointer to the logger
	name       string      // Subsystem name
	endpoints  []*endpoint.Endpoint
}

// register registers the routes for the homepage
func (hps *HomePageSubsystem) register() {
	if !hps.Enabled() {
		hps.Logger.Info("subsystem is disabled", zap.String("subsystem", hps.Name()))
		return
	}

	// Routes
	base := ""
	secure := hps.Controller.Config.Auth.DefaultSecure
	hps.endpoints = []*endpoint.Endpoint{
		{Path: fmt.Sprintf("%s/", base), Action: endpoint.ActionRead, Function: hps.homeHandler, UsesAuth: secure, ExpectedArgs: nil, ExpectedBody: nil, ExpectedOutput: HomeOutput{}},
	}
}

// Enabled returns whether the homepage subsystem is enabled
func (hps *HomePageSubsystem) Enabled() bool {
	return true
}

// Name returns the name of the homepage subsystem
func (hps *HomePageSubsystem) Name() string {
	return hps.name
}

// Endpoints returns an array of endpoints that this Subsystem handles
func (hps *HomePageSubsystem) Endpoints() []*endpoint.Endpoint {
	return hps.endpoints
}

func (hps *HomePageSubsystem) homeHandler(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		response := HomeOutput{
			Message:   fmt.Sprintf("welcome to the %s", hps.Controller.Config.Internal.ProductName),
			Version:   version.Current().String(),
			Endpoints: hps.Controller.EndpointList.Endpoints,
		}
		return response, nil
	}, "dtac-agentd information")
}
