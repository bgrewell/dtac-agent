package basic

import (
	"encoding/json"
	"fmt"
	"github.com/bgrewell/dtac-agent/internal/controller"
	"github.com/bgrewell/dtac-agent/internal/helpers"
	"github.com/bgrewell/dtac-agent/internal/interfaces"
	"github.com/bgrewell/dtac-agent/internal/version"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"go.uber.org/zap"
)

// HomeOutput is the struct for the home page output
type HomeOutput struct {
	Message   string                 `json:"message"`
	Version   string                 `json:"version"`
	Endpoints map[string]interface{} `json:"visible_endpoints"`
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
	authz := endpoint.AuthGroupGuest.String()
	hps.endpoints = []*endpoint.Endpoint{
		endpoint.NewEndpoint(fmt.Sprintf("%s/", base), endpoint.ActionRead, "general dtac agent information", hps.homeHandler, secure, authz, endpoint.WithOutput(HomeOutput{})),
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

func (hps *HomePageSubsystem) homeHandler(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		visibleEndpoints := hps.Controller.EndpointList.GetVisibleEndpoints(in)
		response := HomeOutput{
			Message: fmt.Sprintf("welcome to the %s", hps.Controller.Config.Internal.ProductName),
			Version: version.Current().String(),
			Endpoints: map[string]interface{}{
				"count":     len(visibleEndpoints),
				"endpoints": visibleEndpoints,
			},
		}
		return json.Marshal(response)
	}, "dtac-agentd information")
}
