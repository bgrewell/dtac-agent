package plugins

import (
	"errors"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"net/rpc/jsonrpc"
	"path"
	"strings"
)

type LoadUnloadArgs struct {
	Name string `json:"name"`
}

type DefaultPluginLoader struct {
	PluginDirectory         string
	PluginConfigs           map[string]*PluginConfig
	loadUnconfiguredPlugins bool
	plugins                 map[string]*PluginInfo
	routeMap                map[string]*HandlerEntry
	endpoints               []*endpoint.Endpoint
	pluginRoot              string
}

func (pl *DefaultPluginLoader) Initialize(secure bool) (loadedPlugins []*PluginInfo, err error) {

	// List plugins
	loadedPlugins = make([]*PluginInfo, 0)
	plugs, err := pl.ListPlugins()
	if err != nil {
		return nil, err
	}

	for _, plug := range plugs {

		if config, exists := pl.PluginConfigs[plug]; exists && config.Enabled || pl.loadUnconfiguredPlugins {

			// Launch plugins
			info, err := pl.LaunchPlugin(pl.PluginConfigs[plug])
			if err != nil {
				return nil, err
			}

			// Register plugins
			err = pl.RegisterPlugin(info.Name)
			if err != nil {
				return nil, err
			}
			loadedPlugins = append(loadedPlugins, info)
		}
	}

	// TODO: Need a better way to handle the secure flag
	// Register control routes GET methods are just there for ease of use
	endpoints := []*endpoint.Endpoint{
		{Path: "load", Action: endpoint.ActionRead, UsesAuth: secure, Function: pl.Load, ExpectedArgs: LoadUnloadArgs{}, ExpectedBody: nil, ExpectedOutput: nil},
		{Path: "load", Action: endpoint.ActionCreate, UsesAuth: secure, Function: pl.Load, ExpectedArgs: LoadUnloadArgs{}, ExpectedBody: nil, ExpectedOutput: nil},
		{Path: "unload", Action: endpoint.ActionRead, UsesAuth: secure, Function: pl.Unload, ExpectedArgs: LoadUnloadArgs{}, ExpectedBody: nil, ExpectedOutput: nil},
		{Path: "unload", Action: endpoint.ActionCreate, UsesAuth: secure, Function: pl.Unload, ExpectedArgs: LoadUnloadArgs{}, ExpectedBody: nil, ExpectedOutput: nil},
	}
	pl.endpoints = append(pl.endpoints, endpoints...)

	return loadedPlugins, nil
}

// ListPlugins returns a list of all plugins in the plugin directory
func (pl *DefaultPluginLoader) ListPlugins() (plugins []string, err error) {
	return utility.FindPlugins(pl.PluginDirectory, "*.plugin")
}

// LaunchPlugin launches a plugin and returns the info on the running plugin
func (pl *DefaultPluginLoader) LaunchPlugin(config *PluginConfig) (info *PluginInfo, err error) {
	info, err = executePlugin(config)
	if err != nil {
		return nil, err
	}
	pl.plugins[info.Name] = info
	return info, err
}

// RegisterPlugin registers the plugin routes with Gin
func (pl *DefaultPluginLoader) RegisterPlugin(pluginName string) (err error) {

	if plug, ok := pl.plugins[pluginName]; !ok {
		return errors.New(fmt.Sprintf("no plugin was found with the name: %s", pluginName))
	} else {
		// Connect the rpc client
		plug.Rpc, err = jsonrpc.Dial(plug.Proto, fmt.Sprintf("%s:%d", plug.Ip, plug.Port))
		if err != nil {
			return err
		}
		// Register the plugin
		ra := RegisterArgs{
			Config: pl.plugins[pluginName].PluginConfig.Config,
		}
		rr := &RegisterReply{}
		err = plug.Rpc.Call(fmt.Sprintf("%s.Register", plug.Name), ra, rr)
		if err != nil {
			return err
		}

		// Record routes
		plug.Endpoints = rr.Endpoints
		endpoints := make([]*endpoint.Endpoint, 0)
		for _, ep := range plug.Endpoints {
			// Register endpoints
			ep.Path = path.Join(plug.RootPath, ep.Path)
			endpoints = append(endpoints, ep.Endpoint)

			// Record route map
			key := fmt.Sprintf("%s:%s", ep.Action, ep.Path)
			entry := &HandlerEntry{
				PluginName: plug.Name,
				HandleFunc: ep.FunctionName,
			}
			pl.routeMap[key] = entry
		}
		pl.endpoints = append(pl.endpoints, endpoints...)
	}

	return nil
}

// UnregisterPlugin is used to unregister the plugin from Gin
func (pl *DefaultPluginLoader) UnregisterPlugin(pluginName string) (err error) {
	if _, ok := pl.plugins[pluginName]; !ok {
		return errors.New("plugin not found")
	} else {
		return nil
	}
}

// ClosePlugin is used to stop the plugin process
func (pl *DefaultPluginLoader) ClosePlugin(pluginName string) (err error) {
	if plug, ok := pl.plugins[pluginName]; !ok {
		return errors.New("plugin not found")
	} else {
		token := *plug.CancelToken
		token()
	}
	return nil
}

// Load is used to load a plugin by name
func (pl *DefaultPluginLoader) Load(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		if m := in.Params["name"]; m[0] != "" {
			name := m[0]
			if plug, ok := pl.plugins[name]; ok {
				if !plug.HasExited {
					return nil, errors.New("plugin is already loaded")
				}
				_, err = pl.LaunchPlugin(pl.plugins[name].PluginConfig)
				if err != nil {
					return nil, fmt.Errorf("failed to launch plugin: %s", err)
				}
				err = pl.RegisterPlugin(name)
				if err != nil {
					return nil, fmt.Errorf("failed to register plugin: %s", err)
				}
			} else {
				return nil, fmt.Errorf("no plugin with the name %s found", name)
			}

			return "plugin loaded", nil
		}

		return nil, errors.New("missing 'name' parameter specifying the plugin name")

	}, "plugin loaded")
}

// Unload is used to unload a plugin by name
func (pl *DefaultPluginLoader) Unload(in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {
	return helpers.HandleWrapper(in, func() (interface{}, error) {
		if m := in.Params["name"]; m[0] != "" {
			name := m[0]
			if plug, ok := pl.plugins[name]; ok {
				if plug.HasExited {
					return nil, errors.New("plugin is already unloaded")
				}
				err = pl.UnregisterPlugin(name)
				if err != nil {
					return nil, fmt.Errorf("failed to unregister plugin: %s", err)
				}
				err = pl.ClosePlugin(name)
				if err != nil {
					return nil, fmt.Errorf("failed to close plugin: %s", err)
				}
			} else {
				return nil, fmt.Errorf("no plugin with the name %s found", name)
			}

			return "plugin unloaded", nil
		}

		return nil, errors.New("missing 'name' parameter specifying the plugin name")

	}, "plugin unloaded")
}

// Endpoints returns a list of all the endpoints that are registered with the plugin loader
func (pl *DefaultPluginLoader) Endpoints() []*endpoint.Endpoint {
	return pl.endpoints
}

// CallShim is used to make a call into a plugins function. It acts as a shim between the main internal API and the
// plugin.
func (pl *DefaultPluginLoader) CallShim(ep *endpoint.Endpoint, in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error) {

	// Extract the RouteKey
	keyPath := strings.TrimLeft(strings.Replace(ep.Path, pl.pluginRoot, "", 1), "/")
	routeKey := fmt.Sprintf("%s:%s", ep.Action, keyPath)

	// Get the HandlerEntry
	if handler, ok := pl.routeMap[routeKey]; !ok {
		return nil, errors.New("a handler was not found for the requested resource")
	} else {
		// Make sure plugin isn't canceled
		if pl.plugins[handler.PluginName].HasExited {
			return nil, fmt.Errorf("the plugin has exited with code: %d", pl.plugins[handler.PluginName].ExitCode)
		}

		// Get the plugin
		plug := pl.plugins[handler.PluginName]

		// Clear the context from the input object as it can't be carried across an RPC call
		storedCtx := in.Context
		in.Context = nil

		// Make the rpc call pass in *endpoint.InputArgs and a *endpoint.ReturnVal for the reply
		err = plug.Rpc.Call(fmt.Sprintf("%s.%s", plug.Name, handler.HandleFunc), in, &out)
		if err != nil {
			return nil, fmt.Errorf("failed to call plugin function: %s", err)
		}
		out.Context = storedCtx

		return out, nil
	}
}
