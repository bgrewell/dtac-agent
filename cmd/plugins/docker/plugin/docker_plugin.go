package plugin

import (
	"encoding/json"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/docker/plugin/internal"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"github.com/shirou/gopsutil/docker"
	_ "net/http/pprof" // Used for remote debugging of the plugin
	"reflect"
	"strconv"
)

// DockerIDMessage is just a simple helper struct to encapsulate the id parameter for docker function validation
type DockerIDMessage struct {
	ID []string `json:"id"`
}

// This sets a non-existent variable to the interface type of plugin then attempts to assign
// a pointer to DockerPlugin to it. This isn't needed, but it's a good way to ensure that the
// DockerPlugin struct implements the Plugin interface. If there are missing functions, this
// will fail to compile.
var _ plugins.Plugin = &DockerPlugin{}

// NewDockerPlugin is a constructor that returns a new instance of the DockerPlugin
func NewDockerPlugin() *DockerPlugin {
	// Uncommenting the following anonymous function will allow remote debugging of the plugin by attaching a debugger
	// like the one built into goland. This is useful for debugging plugins.
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	// Create a new instance of the plugin
	hp := &DockerPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}

	// Create a new client wrapper
	client, err := internal.NewDockerClientWrapper()
	if err != nil {
		hp.Log(plugins.LevelFatal, "failed to get docker client", map[string]string{"error": err.Error()})
		return nil
	}
	hp.client = client

	// Return the new instance
	return hp
}

// DockerPlugin is the plugin struct that implements the Plugin interface
type DockerPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	client *internal.DockerClientWrapper
}

// Name returns the name of the plugin type
// NOTE: this is intentionally not a pointer receiver otherwise it wouldn't work. This must be set at your plugin struct
// level. otherwise it will return the type of the PluginBase struct instead.
func (p DockerPlugin) Name() string {
	t := reflect.TypeOf(p)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (p *DockerPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Convert the config json to a map. If you have a specific configuration type you should unmarshal into that type
	var config map[string]interface{}
	err := json.Unmarshal([]byte(request.Config), &config)
	if err != nil {
		return err
	}

	// Declare our endpoint(s)
	authz := endpoint.AuthGroupAdmin.String()
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("/", endpoint.ActionRead, "this endpoint returns docker stats", p.Stats, request.DefaultSecure, authz, endpoint.WithOutput(&[]docker.CgroupDockerStat{})),
		endpoint.NewEndpoint("/images", endpoint.ActionRead, "this endpoint returns docker images", p.ListImages, request.DefaultSecure, authz, endpoint.WithOutput(&[]internal.ImageInfo{})),
		endpoint.NewEndpoint("/configs", endpoint.ActionRead, "this endpoint returns docker config", p.ListConfigs, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/containers", endpoint.ActionRead, "this endpoint returns docker containers", p.ListContainers, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/nodes", endpoint.ActionRead, "this endpoint returns docker nodes", p.ListNodes, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/networks", endpoint.ActionRead, "this endpoint returns docker networks", p.ListNetworks, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/plugins", endpoint.ActionRead, "this endpoint returns docker plugins", p.ListPlugins, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/secrets", endpoint.ActionRead, "this endpoint returns docker secrets", p.ListSecrets, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/services", endpoint.ActionRead, "this endpoint returns docker services", p.ListServices, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/tasks", endpoint.ActionRead, "this endpoint returns docker tasks", p.ListTasks, request.DefaultSecure, authz),
		endpoint.NewEndpoint("/volumes", endpoint.ActionRead, "this endpoint returns docker volumes", p.ListVolumes, request.DefaultSecure, authz),
	}

	// Register them with the plugin
	p.RegisterMethods(endpoints)

	// Convert to plugin endpoints and return
	for _, ep := range endpoints {
		aep := utility.ConvertEndpointToPluginEndpoint(ep)
		reply.Endpoints = append(reply.Endpoints, aep)
	}

	// Print out a log message
	p.Log(plugins.LevelInfo, "docker plugin registered", map[string]string{"endpoint_count": strconv.Itoa(len(endpoints))})

	// Return no error
	return nil
}

// ListImages is the handler for the docker images route
func (p *DockerPlugin) ListImages(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {

		options, err := internal.ParseListImageOptions(in.Parameters)
		if err != nil {
			p.Log(plugins.LevelError, "failed to parse image list options", map[string]string{"error": err.Error()})
		}
		images, err := p.client.ListImages(options...)
		if err != nil {
			return nil, err
		}
		return json.Marshal(images)
	}, "docker image information")
}

// ListConfigs is the handler for the docker list configs route
func (p *DockerPlugin) ListConfigs(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {

		//options, err := internal.ParseListImageOptions(in.Parameters)
		configs, err := p.client.ListConfigs()
		if err != nil {
			return nil, err
		}
		return json.Marshal(configs)
	}, "docker config information")
}

// ListContainers is the handler for the docker list containers route
func (p *DockerPlugin) ListContainers(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {

		options, err := internal.ParseListContainerOptions(in.Parameters)
		if err != nil {
			p.Log(plugins.LevelError, "failed to parse container list options", map[string]string{"error": err.Error()})
		}
		containers, err := p.client.ListContainers(options...)
		if err != nil {
			return nil, err
		}
		return json.Marshal(containers)
	}, "docker container information")
}

// ListNodes is the handler for the docker list nodes route
func (p *DockerPlugin) ListNodes(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		nodes, err := p.client.ListNodes()
		if err != nil {
			return nil, err
		}
		return json.Marshal(nodes)
	}, "docker node information")
}

// ListNetworks is the handler for the docker list networks route
func (p *DockerPlugin) ListNetworks(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		networks, err := p.client.ListNetworks()
		if err != nil {
			return nil, err
		}
		return json.Marshal(networks)
	}, "docker network information")
}

// ListPlugins is the handler for the docker list plugins route
func (p *DockerPlugin) ListPlugins(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		plugins, err := p.client.ListPlugins()
		if err != nil {
			return nil, err
		}
		return json.Marshal(plugins)
	}, "docker plugins information")
}

// ListSecrets is the handler for the docker list secrets route
func (p *DockerPlugin) ListSecrets(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		secrets, err := p.client.ListSecrets()
		if err != nil {
			return nil, err
		}
		return json.Marshal(secrets)
	}, "docker secrets information")
}

// ListServices is the handler for the docker list services route
func (p *DockerPlugin) ListServices(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		services, err := p.client.ListServices()
		if err != nil {
			return nil, err
		}
		return json.Marshal(services)
	}, "docker services information")
}

// ListTasks is the handler for the docker list tasks route
func (p *DockerPlugin) ListTasks(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		tasks, err := p.client.ListTasks()
		if err != nil {
			return nil, err
		}
		return json.Marshal(tasks)
	}, "docker tasks information")
}

// ListVolumes is the handler for the docker list volumes route
func (p *DockerPlugin) ListVolumes(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapper(in, func() ([]byte, error) {
		volumes, err := p.client.ListVolumes()
		if err != nil {
			return nil, err
		}
		return json.Marshal(volumes)
	}, "docker volumes information")
}

// Stats is the handler for the docker stats route
func (p *DockerPlugin) Stats(in *endpoint.Request) (out *endpoint.Response, err error) {
	return p.ListContainers(in)
}
