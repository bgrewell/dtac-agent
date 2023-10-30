package endpoints

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"go.uber.org/zap"
)

// NewEndpointList creates a new instance of the RouteList struct
func NewEndpointList(cfg *config.Configuration, log *zap.Logger) *EndpointList {
	httpList := EndpointList{
		Config: cfg,
		Logger: log.With(zap.String("module", "route_list")),
	}
	return &httpList
}

// EndpointList is the struct for the api endpoint list
type EndpointList struct {
	Endpoints []*endpoint.Endpoint  `json:"endpoints" yaml:"endpoints"`
	Config    *config.Configuration `json:"-" yaml:"-"`
	Logger    *zap.Logger           `json:"-" yaml:"-"`
}

// AddEndpoints inserts new endpoints into the endpoint list
func (el *EndpointList) AddEndpoints(endpoints []*endpoint.Endpoint) {
	el.Endpoints = append(el.Endpoints, endpoints...)
}
