package endpoints

import (
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"go.uber.org/zap"
	"strings"
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

// GetVisibleEndpoints returns a list of endpoints that are visible to the user
func (el *EndpointList) GetVisibleEndpoints(in *endpoint.Request) (visibleEndpoints []*endpoint.Endpoint) {
	visibleEndpoints = make([]*endpoint.Endpoint, 0)
	roleMap := make(map[string]bool)
	if roles, ok := in.Metadata[types.ContextAuthRoles.String()]; ok {
		for _, role := range strings.Split(roles, ",") {
			switch role {
			case endpoint.AuthGroupAdmin.String():
				roleMap = map[string]bool{"admin": true, "operator": true, "user": true, "guest": true}
			case endpoint.AuthGroupOperator.String():
				roleMap = map[string]bool{"operator": true, "user": true, "guest": true}
			case endpoint.AuthGroupUser.String():
				roleMap = map[string]bool{"user": true, "guest": true}
			case endpoint.AuthGroupGuest.String():
				roleMap = map[string]bool{"guest": true}
			default:
				roleMap = map[string]bool{}
			}
		}
	}

	for _, ep := range el.Endpoints {
		if _, hasAccess := roleMap[ep.AuthGroup]; !ep.Secure || hasAccess {
			visibleEndpoints = append(visibleEndpoints, ep)
		}
	}

	return visibleEndpoints
}
