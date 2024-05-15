package plugin

import (
	"encoding/json"
	"errors"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"reflect"
	"strconv"
)

type Something struct {
}

// Ensure DSCPlugin implements the Plugin interface
var _ plugins.Plugin = &DSCPlugin{}

// NewDSCPlugin is a constructor that returns a new instance of the DSCPlugin
func NewDSCPlugin() *DSCPlugin {
	dp := &DSCPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}
	dp.SetRootPath("dsc")
	return dp
}

// DSCPlugin is the plugin struct that implements the Plugin interface
type DSCPlugin struct {
	plugins.PluginBase
}

// Name returns the name of the plugin type
func (d DSCPlugin) Name() string {
	t := reflect.TypeOf(d)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (d *DSCPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Convert the config json to a map. If you have a specific configuration type you should unmarshal into that type
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return err
	}

	authz := endpoint.AuthGroupAdmin.String()
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("adusers", endpoint.ActionRead, "this endpoint returns a list of active directory users", d.GetADUsers, request.DefaultSecure, authz, endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionRead, "this endpoint returns an active directory user", d.GetADUser, request.DefaultSecure, authz, endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionCreate, "this endpoint creates a new active directory user", d.CreateADUser, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionWrite, "this endpoint updates an active directory user", d.UpdateADUser, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionDelete, "this endpoint deletes an active directory user", d.DeleteADUser, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adgroups", endpoint.ActionRead, "this endpoint returns a list of active directory groups", d.GetADGroups, request.DefaultSecure, authz, endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adgroup", endpoint.ActionCreate, "this endpoint creates a new active directory group", d.CreateADGroup, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adgroup", endpoint.ActionWrite, "this endpoint updates an active directory group", d.UpdateADGroup, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adgroup", endpoint.ActionDelete, "this endpoint deletes an active directory group", d.DeleteADGroup, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adous", endpoint.ActionRead, "this endpoint returns a list of active directory organizational units", d.GetADOrganizationalUnits, request.DefaultSecure, authz, endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adou", endpoint.ActionCreate, "this endpoint creates a new active directory organizational unit", d.CreateADOrganizationalUnit, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adou", endpoint.ActionWrite, "this endpoint updates an active directory organizational unit", d.UpdateADOrganizationalUnit, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
		endpoint.NewEndpoint("adou", endpoint.ActionDelete, "this endpoint deletes an active directory organizational unit", d.DeleteADOrganizationalUnit, request.DefaultSecure, authz, endpoint.WithBody(&Something{}), endpoint.WithOutput(&Something{})),
	}

	d.RegisterMethods(endpoints)

	for _, ep := range endpoints {
		aep := utility.ConvertEndpointToPluginEndpoint(ep)
		reply.Endpoints = append(reply.Endpoints, aep)
	}

	d.Log(plugins.LevelInfo, "dsc plugin registered", map[string]string{"endpoint_count": strconv.Itoa(len(endpoints))})
	return nil
}

func (d DSCPlugin) GetADUsers(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) GetADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) CreateADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) UpdateADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) DeleteADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) CreateADGroup(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) UpdateADGroup(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) DeleteADGroup(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) CreateADOrganizationalUnit(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) UpdateADOrganizationalUnit(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) DeleteADOrganizationalUnit(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) GetADGroups(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DSCPlugin) GetADOrganizationalUnits(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}
