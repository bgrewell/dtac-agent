package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	exec "github.com/bgrewell/go-execute/v2"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"github.com/bgrewell/dtac-agent/pkg/plugins/utility"
	"reflect"
	"strconv"
)

// Ensure DomainControllerPlugin implements the Plugin interface
var _ plugins.Plugin = &DomainControllerPlugin{}

// Expose an executor for the plugin to use
var executor = exec.NewExecutor()

// NewDomainControllerPlugin is a constructor that returns a new instance of the DomainControllerPlugin
func NewDomainControllerPlugin() *DomainControllerPlugin {
	dp := &DomainControllerPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
		headers: map[string][]string{
			"X-PLUGIN-NAME": {"active-directory"},
		},
	}
	dp.SetRootPath("dsc")
	return dp
}

// DomainControllerPlugin is the plugin struct that implements the Plugin interface
type DomainControllerPlugin struct {
	plugins.PluginBase
	headers map[string][]string
}

// Name returns the name of the plugin type
func (d DomainControllerPlugin) Name() string {
	t := reflect.TypeOf(d)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (d *DomainControllerPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Convert the config json to a map. If you have a specific configuration type you should unmarshal into that type
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return err
	}

	authz := endpoint.AuthGroupAdmin.String()
	endpoints := []*endpoint.Endpoint{
		// Users
		endpoint.NewEndpoint("adusers", endpoint.ActionRead, "this endpoint returns a list of active directory users", d.GetADUsers, request.DefaultSecure, authz, endpoint.WithOutput(&[]ADUser{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionRead, "this endpoint returns an active directory user", d.GetADUser, request.DefaultSecure, authz, endpoint.WithParameters(&ADUserParam{}), endpoint.WithOutput(&ADUser{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionCreate, "this endpoint creates a new active directory user", d.CreateADUser, request.DefaultSecure, authz, endpoint.WithBody(&NewADUserObj{}), endpoint.WithOutput(&ADUser{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionWrite, "this endpoint updates an active directory user", d.UpdateADUser, request.DefaultSecure, authz, endpoint.WithBody(&UpdateADUserObj{}), endpoint.WithOutput(&ADUser{})),
		endpoint.NewEndpoint("aduser", endpoint.ActionDelete, "this endpoint deletes an active directory user", d.DeleteADUser, request.DefaultSecure, authz, endpoint.WithParameters(&ADUserParam{}), endpoint.WithOutput(&ADUser{})),
		endpoint.NewEndpoint("aduser/sshkeys", endpoint.ActionRead, "this endpoint returns a list of ssh keys for an active directory user", d.GetADUserSSHKeys, request.DefaultSecure, authz, endpoint.WithParameters(&ADUserParam{}), endpoint.WithOutput(&[]string{})),
		endpoint.NewEndpoint("aduser/sshkey", endpoint.ActionCreate, "this endpoint adds an ssh key to an active directory user", d.AddADUserSSHKey, request.DefaultSecure, authz, endpoint.WithBody(&ADUserSSHKey{}), endpoint.WithOutput(&[]string{})),
		endpoint.NewEndpoint("aduser/sshkey", endpoint.ActionDelete, "this endpoint deletes an ssh key from an active directory user", d.DeleteADUserSSHKey, request.DefaultSecure, authz, endpoint.WithBody(&ADUserSSHKey{}), endpoint.WithOutput(&[]string{})),
		endpoint.NewEndpoint("aduser/importsshkey", endpoint.ActionCreate, "this endpoint imports an ssh key to an active directory user", d.ImportADUserSSHKey, request.DefaultSecure, authz, endpoint.WithBody(&ADUserSSHKeyImport{}), endpoint.WithOutput(&[]string{})),

		// Groups
		endpoint.NewEndpoint("adgroups", endpoint.ActionRead, "this endpoint returns a list of active directory groups", d.GetADGroups, request.DefaultSecure, authz, endpoint.WithOutput(&[]ADGroup{})),
		endpoint.NewEndpoint("adgroup", endpoint.ActionRead, "this endpoint returns an active directory group", d.GetADGroup, request.DefaultSecure, authz, endpoint.WithParameters(&ADGroupName{}), endpoint.WithOutput(&ADGroup{})),
		endpoint.NewEndpoint("adgroup", endpoint.ActionCreate, "this endpoint creates a new active directory group", d.CreateADGroup, request.DefaultSecure, authz, endpoint.WithBody(&ADGroup{}), endpoint.WithOutput(&ADGroup{})),
		endpoint.NewEndpoint("adgroup", endpoint.ActionWrite, "this endpoint updates an active directory group", d.UpdateADGroup, request.DefaultSecure, authz, endpoint.WithBody(&ADGroup{}), endpoint.WithOutput(&ADGroup{})),
		endpoint.NewEndpoint("adgroup", endpoint.ActionDelete, "this endpoint deletes an active directory group", d.DeleteADGroup, request.DefaultSecure, authz, endpoint.WithBody(&DeleteADGroup{}), endpoint.WithOutput(&ADGroup{})),
		// Organizational Units
		endpoint.NewEndpoint("adous", endpoint.ActionRead, "this endpoint returns a list of active directory organizational units", d.GetADOrganizationalUnits, request.DefaultSecure, authz, endpoint.WithOutput(&[]ADOrganizationalUnit{})),
		endpoint.NewEndpoint("adou", endpoint.ActionRead, "this endpoint returns an active directory organizational unit", d.GetADOrganizationalUnit, request.DefaultSecure, authz, endpoint.WithParameters(&ADOrganizationalUnitName{}), endpoint.WithOutput(&ADOrganizationalUnit{})),
		endpoint.NewEndpoint("adou", endpoint.ActionCreate, "this endpoint creates a new active directory organizational unit", d.CreateADOrganizationalUnit, request.DefaultSecure, authz, endpoint.WithBody(&NewADOrganizationalUnit{}), endpoint.WithOutput(&ADOrganizationalUnit{})),
		endpoint.NewEndpoint("adou", endpoint.ActionWrite, "this endpoint updates an active directory organizational unit", d.UpdateADOrganizationalUnit, request.DefaultSecure, authz, endpoint.WithBody(&ADOrganizationalUnit{}), endpoint.WithOutput(&ADOrganizationalUnit{})),
		endpoint.NewEndpoint("adou", endpoint.ActionDelete, "this endpoint deletes an active directory organizational unit", d.DeleteADOrganizationalUnit, request.DefaultSecure, authz, endpoint.WithBody(&ADOrganizationalUnitName{}), endpoint.WithOutput(&ADOrganizationalUnit{})),
		// Sudo Role
		endpoint.NewEndpoint("adsudos", endpoint.ActionRead, "this endpoint returns a list of active directory sudoers", d.GetADSudoRoles, request.DefaultSecure, authz, endpoint.WithOutput(&[]ADSudoRole{})),
		endpoint.NewEndpoint("adsudo", endpoint.ActionRead, "this endpoint returns a directory sudoer entry", d.GetADSudoRole, request.DefaultSecure, authz, endpoint.WithParameters(&ADSudoRoleName{}), endpoint.WithOutput(&ADSudoRole{})),
		endpoint.NewEndpoint("adsudo", endpoint.ActionCreate, "this endpoint creates a new active directory sudoer", d.CreateADSudoRole, request.DefaultSecure, authz, endpoint.WithBody(&NewSudoRole{}), endpoint.WithOutput(&ADSudoRole{})),
		endpoint.NewEndpoint("adsudo", endpoint.ActionWrite, "this endpoint updates an active directory sudoer", d.UpdateADSudoRole, request.DefaultSecure, authz, endpoint.WithBody(&ADSudoRole{}), endpoint.WithOutput(&ADSudoRole{})),
		endpoint.NewEndpoint("adsudo", endpoint.ActionDelete, "this endpoint deletes an active directory sudoer", d.DeleteADSudoRole, request.DefaultSecure, authz, endpoint.WithBody(&DeleteSudoRole{}), endpoint.WithOutput(&ADSudoRole{})),
	}

	d.RegisterMethods(endpoints)

	for _, ep := range endpoints {
		aep := utility.ConvertEndpointToPluginEndpoint(ep)
		reply.Endpoints = append(reply.Endpoints, aep)
	}

	d.Log(plugins.LevelInfo, "dsc plugin registered", map[string]string{"endpoint_count": strconv.Itoa(len(endpoints))})
	return nil
}

func (d DomainControllerPlugin) GetADUsers(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		users, err := ListADUsers()
		if err != nil {
			return nil, nil, err
		}
		out, err := json.Marshal(users)
		if err != nil {
			return nil, nil, err
		}
		return d.headers, out, nil
	}, "list of active directory users")
}

func (d DomainControllerPlugin) GetADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		if username, ok := in.Parameters["username"]; ok {
			user, err := ListADUser(username[0])
			if err != nil {
				return nil, nil, err
			}
			out, err := json.Marshal(user)
			if err != nil {
				return nil, nil, err
			}
			return d.headers, out, nil
		}

		return nil, nil, fmt.Errorf("the parameter username is required")

	}, "list of active directory user details")
}

func (d DomainControllerPlugin) CreateADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		user := &NewADUserObj{}
		if in.Body != nil {
			err := json.Unmarshal(in.Body, user)
			if err != nil {
				return nil, nil, err
			}
		}

		err := CreateADUser(user)
		if err != nil {
			return nil, nil, err
		}

		verifyUser, err := ListADUser(user.Username)
		if err != nil {
			return nil, nil, err
		}

		out, err := json.Marshal(verifyUser)
		if err != nil {
			return nil, nil, err
		}

		return d.headers, out, nil
	}, "create a new active directory user")
}

func (d DomainControllerPlugin) UpdateADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		user := &NewADUserObj{}
		if in.Body != nil {
			err := json.Unmarshal(in.Body, user)
			if err != nil {
				return nil, nil, err
			}
		}

		err := UpdateADUser(user)
		if err != nil {
			return nil, nil, err
		}

		verifyUser, err := ListADUser(user.Username)
		if err != nil {
			return nil, nil, err
		}

		out, err := json.Marshal(verifyUser)
		if err != nil {
			return nil, nil, err
		}

		return d.headers, out, nil
	}, "update an active directory user")
}

func (d DomainControllerPlugin) DeleteADUser(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		if username, ok := in.Parameters["username"]; ok {
			err := DeleteADUser(username[0])
			if err != nil {
				return nil, nil, err
			}
			out, err := json.Marshal(map[string]string{
				"username": username[0],
				"status":   "deleted",
			})
			if err != nil {
				return nil, nil, err
			}
			return d.headers, out, nil
		}

		return nil, nil, fmt.Errorf("the parameter username is required")

	}, "delete an active directory user")
}

func (d DomainControllerPlugin) CreateADGroup(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) UpdateADGroup(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) DeleteADGroup(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) CreateADOrganizationalUnit(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		ou := &NewADOrganizationalUnit{}
		if in.Body != nil {
			err := json.Unmarshal(in.Body, ou)
			if err != nil {
				return nil, nil, err
			}
		}

		unit, err := CreateADOrganizationalUnit(ou)
		if err != nil {
			return nil, nil, err
		}

		out, err := json.Marshal(unit)
		if err != nil {
			return nil, nil, err
		}

		return d.headers, out, nil
	}, "create a new active directory organizational unit")
}

func (d DomainControllerPlugin) UpdateADOrganizationalUnit(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) DeleteADOrganizationalUnit(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		if name, ok := in.Parameters["name"]; ok {
			err := DeleteADOrganizationalUnit(name[0])
			if err != nil {
				return nil, nil, err
			}
			out, err := json.Marshal(map[string]string{
				"name":   name[0],
				"status": "deleted",
			})
			if err != nil {
				return nil, nil, err
			}
			return d.headers, out, nil
		}

		return nil, nil, fmt.Errorf("the parameter name is required")

	}, "list of active directory organizational unit details")
}

func (d DomainControllerPlugin) GetADGroups(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		groups, err := ListADGroups()
		if err != nil {
			return nil, nil, err
		}
		out, err := json.Marshal(groups)
		if err != nil {
			return nil, nil, err
		}
		return d.headers, out, nil
	}, "list of active directory groups")
}

func (d DomainControllerPlugin) GetADOrganizationalUnits(in *endpoint.Request) (out *endpoint.Response, err error) {

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		ous, err := ListADOrganizationalUnits()
		if err != nil {
			return nil, nil, err
		}
		out, err := json.Marshal(ous)
		if err != nil {
			return nil, nil, err
		}
		return d.headers, out, nil
	}, "list of active directory organizational units")
}

func (d DomainControllerPlugin) GetADSudoRoles(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		sudoers, err := ListADSudoRoles()
		if err != nil {
			return nil, nil, err
		}
		out, err := json.Marshal(sudoers)
		if err != nil {
			return nil, nil, err
		}
		return d.headers, out, nil
	}, "list of active directory sudo roles")
}

func (d DomainControllerPlugin) CreateADSudoRole(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) UpdateADSudoRole(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) DeleteADSudoRole(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) GetADSudoRole(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		if name, ok := in.Parameters["name"]; ok {
			user, err := ListADSudoRole(name[0])
			if err != nil {
				return nil, nil, err
			}
			out, err := json.Marshal(user)
			if err != nil {
				return nil, nil, err
			}
			return d.headers, out, nil
		}

		return nil, nil, fmt.Errorf("the parameter name is required")

	}, "list of active directory sudorole details")
}

func (d DomainControllerPlugin) GetADGroup(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		if name, ok := in.Parameters["name"]; ok {
			group, err := ListADGroup(name[0])
			if err != nil {
				return nil, nil, err
			}
			out, err := json.Marshal(group)
			if err != nil {
				return nil, nil, err
			}
			return d.headers, out, nil
		}

		return nil, nil, fmt.Errorf("the parameter name is required")

	}, "list of active directory group details")
}

func (d DomainControllerPlugin) GetADOrganizationalUnit(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		if name, ok := in.Parameters["name"]; ok {
			unit, err := ListADOrganizationalUnit(name[0])
			if err != nil {
				return nil, nil, err
			}
			out, err := json.Marshal(unit)
			if err != nil {
				return nil, nil, err
			}
			return d.headers, out, nil
		}

		return nil, nil, fmt.Errorf("the parameter name is required")

	}, "list of active directory organizational unit details")
}

func (d DomainControllerPlugin) GetADUserSSHKeys(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) AddADUserSSHKey(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) DeleteADUserSSHKey(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}

func (d DomainControllerPlugin) ImportADUserSSHKey(in *endpoint.Request) (out *endpoint.Response, err error) {
	return nil, errors.New("Not implemented")
}
