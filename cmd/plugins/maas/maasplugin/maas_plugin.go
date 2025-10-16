package maasplugin

import (
	"encoding/json"
	"fmt"
	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maasplugin/engine"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maasplugin/structs"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"log"
	"os"
	"reflect"
	"strconv"
)

// Ensure that our type meets the requirements for being a plugin
var _ plugins.Plugin = &MAASPlugin{}

func NewMAASPlugin() *MAASPlugin {

	p := &MAASPlugin{
		PluginBase: plugins.PluginBase{
			Methods: make(map[string]endpoint.Func),
		},
	}

	return p
}

// MAASPlugin is the main plugin type
type MAASPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	settings *structs.MAASSettings
	engine   *engine.Engine
}

// Name returns the name of the plugin
func (p MAASPlugin) Name() string {
	t := reflect.TypeOf(p)
	return t.Name()
}

// TODO: Remove later, this was a reminant of the old plugin system
//// RouteRoot returns the root path for the plugin
//func (p MAASPlugin) RouteRoot() string {
//	return "maas"
//}

// Register registers the MAAS plugin with the plugin manager
func (p *MAASPlugin) Register(request *api.RegisterRequest, reply *api.RegisterResponse) error {
	*reply = api.RegisterResponse{Endpoints: make([]*api.PluginEndpoint, 0)}

	// Convert the config json to a map. If you have a specific configuration type you should unmarshal into that type
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(request.Config), &config); err != nil {
		return err
	}

	// Check if the configuration has a logfile and set up logging
	if lf, ok := config["logfile"]; ok {
		if path, ok := lf.(string); ok {
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
			if err != nil {
				p.Log(plugins.LevelError, "failed to open log file", map[string]string{"error": err.Error()})
			} else {
				log.SetOutput(file)
				p.Log(plugins.LevelInfo, "logging to file", map[string]string{"file": path})
			}
		} else {
			p.Log(plugins.LevelError, "logfile must be a string", map[string]string{"type": fmt.Sprintf("%T", lf)})
		}
	}

	// Check if the configuration has a server and set the server address
	var server string
	if sv, ok := config["server"]; ok {
		if s, ok := sv.(string); ok {
			server = s
		} else {
			p.Log(plugins.LevelError, "server must be a string", map[string]string{"type": fmt.Sprintf("%T", sv)})
			server = "localhost"
		}
	} else {
		server = "localhost"
	}

	// Check if the configuration has poll_interval_secs and set the polling interval
	pollSecs := 30
	if pi, ok := config["poll_interval_secs"]; ok {
		if f, ok := pi.(float64); ok {
			pollSecs = int(f)
		} else {
			p.Log(plugins.LevelError, "poll_interval_secs must be a number", map[string]string{"type": fmt.Sprintf("%T", pi)})
		}
	}

	// Check for required plugin configuration keys
	for _, key := range []string{"consumer_token", "auth_token", "auth_signature"} {
		if _, ok := config[key]; !ok {
			return fmt.Errorf("missing required plugin configuration: %s", key)
		}
	}

	// Setup credentials
	p.settings = &structs.MAASSettings{
		Server:          server,
		ConsumerToken:   config["consumer_token"].(string),
		AuthToken:       config["auth_token"].(string),
		AuthSignature:   config["auth_signature"].(string),
		MachinePollSecs: pollSecs,
	}

	// Declare our endpoint(s)
	authz := endpoint.AuthGroupOperator.String()
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("start", endpoint.ActionRead, "Start MAAS polling", p.Start, request.DefaultSecure, authz),
		endpoint.NewEndpoint("stop", endpoint.ActionRead, "Stop MAAS polling", p.Stop, request.DefaultSecure, authz),
		endpoint.NewEndpoint("status", endpoint.ActionRead, "Get polling status", p.Status, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machines", endpoint.ActionRead, "List all machines", p.GetMachines, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machines/ids", endpoint.ActionRead, "List machine IDs", p.GetMachinesIDs, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machines/pools", endpoint.ActionRead, "List machine pools", p.GetMachinesPools, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machines/status", endpoint.ActionRead, "List machine statuses", p.GetMachinesStatuses, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machines/interfaces", endpoint.ActionRead, "List machine interfaces", p.GetMachinesInterfaces, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machine", endpoint.ActionRead, "Get a single machine", p.GetMachine, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machine/id", endpoint.ActionRead, "Get machine by ID", p.GetMachineID, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machine/pool", endpoint.ActionRead, "Get machine pool", p.GetMachinePool, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machine/status", endpoint.ActionRead, "Get machine status", p.GetMachineStatus, request.DefaultSecure, authz),
		endpoint.NewEndpoint("machine/interfaces", endpoint.ActionRead, "Get machine interfaces", p.GetMachineInterfaces, request.DefaultSecure, authz),
		endpoint.NewEndpoint("fabrics", endpoint.ActionRead, "List all fabrics", p.GetFabrics, request.DefaultSecure, authz),
		endpoint.NewEndpoint("fabric", endpoint.ActionRead, "Get a single fabric", p.GetFabric, request.DefaultSecure, authz),
	}

	// Register them with the plugin
	p.RegisterMethods(endpoints)
	for _, ep := range endpoints {
		reply.Endpoints = append(reply.Endpoints, utility.ConvertEndpointToPluginEndpoint(ep))
	}

	// Print out a log message
	p.Log(plugins.LevelInfo, "maas plugin registered", map[string]string{"endpoint_count": strconv.Itoa(len(endpoints))})

	// Auto-start if requested
	if _, ok := config["auto_start"]; ok {
		_, err := p.Start(nil)
		if err != nil {
			return fmt.Errorf("failed to start plugin: %v", err)
		}
	}

	return nil
}

// Start starts the plugin
func (p *MAASPlugin) Start(in *endpoint.Request) (out *endpoint.Response, err error) {

	if p.engine == nil || !p.engine.Running() {
		p.engine = &engine.Engine{
			Settings: p.settings,
		}
	} else {
		return nil, fmt.Errorf("engine is already running. you must stop first")
	}

	err = p.engine.Start()
	if err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}
		msg, err := json.Marshal("{\"status\": \"started\"}")
		if err != nil {
			return nil, nil, err
		}
		return headers, msg, nil
	}, "maas plugin start response")
}

// Stop stops the plugin
func (p *MAASPlugin) Stop(in *endpoint.Request) (out *endpoint.Response, err error) {
	// Clear creds and stop the engine
	if p.engine.Running() {
		err = p.engine.Stop()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("plugin engine is not running")
	}
	return nil, fmt.Errorf("stopped ... but this function is hardcoded for testing and needs to be finished")
}

// Status returns the status of the plugin
func (p *MAASPlugin) Status(in *endpoint.Request) (out *endpoint.Response, err error) {
	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		if p.engine != nil {
			errstr := ""
			if p.engine.ErrDetails() != nil {
				errstr = p.engine.ErrDetails().Error()
			}

			count := 0
			machines := p.engine.Machines()
			if machines != nil {
				count = len(machines)
			}
			status := structs.Status{
				Running:      p.engine.Running(),
				Failed:       p.engine.Failed(),
				ErrDetails:   errstr,
				MachineCount: count,
			}

			o, e := p.Serialize(status)
			if e != nil {
				return nil, nil, e
			}
			return headers, []byte(o), nil
		} else {
			o, e := p.Serialize("plugin has not been started")
			if e != nil {
				return nil, nil, e
			}
			return headers, []byte(o), nil
		}
	}, "maas plugin status message")
}

// CreateMachine asks MAAS to create a new machine using the parameters supplied
// in in.Parameters, then returns the JSON response.
func (p *MAASPlugin) CreateMachine(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		// Pass the incoming form‐fields map straight to the engine
		respBytes, e := p.engine.CreateMachine(in)
		if e != nil {
			return nil, nil, e
		}

		return headers, respBytes, nil
	}, "maas plugin create_machine response")
}

// GetMachines returns a list of machines from the MAAS server
func (p *MAASPlugin) GetMachines(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		machines := p.engine.Machines()
		var value []byte
		if machines == nil {
			v, e := p.Serialize([]string{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			v, e := p.Serialize(machines)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machines response")
}

// GetMachinesIDs returns a list of machine ids from the MAAS server
func (p *MAASPlugin) GetMachinesIDs(in *endpoint.Request) (out *endpoint.Response, err error) {
	// verify plugin is ready
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		// structure for hostname + system ID
		type MachinePlug struct {
			Hostname string `json:"hostname"`
			SystemID string `json:"system_id"`
		}

		machines := p.engine.Machines()
		var value []byte

		// serialize an empty slice if there are no machines
		if machines == nil {
			v, e := p.Serialize([]MachinePlug{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			// build our list of MachinePlug
			ids := make([]MachinePlug, 0, len(machines))
			for _, m := range machines {
				ids = append(ids, MachinePlug{
					Hostname: m.Hostname,
					SystemID: m.SystemID,
				})
			}
			v, e := p.Serialize(ids)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machines_ids response")
}

// GetMachinesPools returns a list of machine pools from the MAAS server
func (p *MAASPlugin) GetMachinesPools(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		type Result struct {
			Machine string      `json:"machine"`
			Pool    interface{} `json:"pool"`
		}

		machines := p.engine.Machines()
		var value []byte

		if machines == nil {
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			results := make([]Result, 0, len(machines))
			for _, m := range machines {
				results = append(results, Result{
					Machine: m.Hostname,
					Pool:    m.Pool,
				})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machines_pools response")
}

// GetMachinesStatuses returns a list of machine statuses from the MAAS server
func (p *MAASPlugin) GetMachinesStatuses(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		type Result struct {
			Machine string      `json:"machine"`
			Status  interface{} `json:"status"`
		}

		machines := p.engine.Machines()
		var value []byte

		if machines == nil {
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			results := make([]Result, 0, len(machines))
			for _, m := range machines {
				results = append(results, Result{
					Machine: m.Hostname,
					Status:  m.Status,
				})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machines_statuses response")
}

// GetMachinesInterfaces returns a list of machine interfaces from the MAAS server
func (p *MAASPlugin) GetMachinesInterfaces(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		type Result struct {
			Machine    string                    `json:"machine"`
			Interfaces []structs.InterfaceStruct `json:"interfaces"`
		}

		machines := p.engine.Machines()
		var value []byte

		if machines == nil {
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			results := make([]Result, 0, len(machines))
			for _, m := range machines {
				results = append(results, Result{
					Machine:    m.Hostname,
					Interfaces: m.InterfaceSet,
				})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machines_interfaces response")
}

// GetMachine returns a machine from the MAAS server
func (p *MAASPlugin) GetMachine(in *endpoint.Request) (out *endpoint.Response, err error) {
	// verify plugin is ready
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		// find matching machines based on request parameters
		matches, err := p.getMatchingMachines(in)
		if err != nil {
			return nil, nil, err
		}

		type Result struct {
			Machine *structs.Machine `json:"machine"`
		}

		var value []byte
		if matches == nil {
			// no matches → empty slice
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			// wrap each Machine in a Result
			results := make([]Result, 0, len(matches))
			for _, m := range matches {
				results = append(results, Result{Machine: m})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machine response")
}

// GetMachineID returns a machine id from the MAAS server
func (p *MAASPlugin) GetMachineID(in *endpoint.Request) (out *endpoint.Response, err error) {
	// verify plugin is ready
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		// find matching machines based on request parameters
		matches, err := p.getMatchingMachines(in)
		if err != nil {
			return nil, nil, err
		}

		type Result struct {
			Machine string `json:"machine"`
			ID      string `json:"id"`
		}

		var value []byte
		if matches == nil {
			// no matches → empty slice
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			// build our list of Machine+ID pairs
			results := make([]Result, 0, len(matches))
			for _, m := range matches {
				results = append(results, Result{
					Machine: m.Hostname,
					ID:      m.SystemID,
				})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machine_id response")
}

// GetMachinePool returns a machine pool from the MAAS server
func (p *MAASPlugin) GetMachinePool(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		// find matching machines based on request parameters
		matches, err := p.getMatchingMachines(in)
		if err != nil {
			return nil, nil, err
		}

		type Result struct {
			Machine string      `json:"machine"`
			Pool    interface{} `json:"pool"`
		}

		var value []byte
		if matches == nil {
			// no matches → empty slice
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			// build our list of Machine+Pool pairs
			results := make([]Result, 0, len(matches))
			for _, m := range matches {
				results = append(results, Result{
					Machine: m.Hostname,
					Pool:    m.Pool,
				})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machine_pool response")
}

// GetMachineStatus returns a machine status from the MAAS server
func (p *MAASPlugin) GetMachineStatus(in *endpoint.Request) (out *endpoint.Response, err error) {
	// verify plugin is ready
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		// find matching machines based on request parameters
		matches, err := p.getMatchingMachines(in)
		if err != nil {
			return nil, nil, err
		}

		type Result struct {
			Machine string      `json:"machine"`
			Status  interface{} `json:"status"`
		}

		var value []byte
		if matches == nil {
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			results := make([]Result, 0, len(matches))
			for _, m := range matches {
				results = append(results, Result{
					Machine: m.Hostname,
					Status:  m.Status,
				})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machine_status response")
}

// GetMachineInterfaces returns a machine interfaces from the MAAS server
func (p *MAASPlugin) GetMachineInterfaces(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		type Result struct {
			Machine    string                    `json:"machine"`
			Interfaces []structs.InterfaceStruct `json:"interfaces"`
		}

		matches, err := p.getMatchingMachines(in)
		if err != nil {
			return nil, nil, err
		}

		var value []byte
		if matches == nil {
			v, e := p.Serialize([]Result{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			results := make([]Result, 0, len(matches))
			for _, m := range matches {
				results = append(results, Result{
					Machine:    m.Hostname,
					Interfaces: m.InterfaceSet,
				})
			}
			v, e := p.Serialize(results)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_machine_interfaces response")
}

// GetFabrics returns a list of fabrics from the MAAS server
func (p *MAASPlugin) GetFabrics(in *endpoint.Request) (out *endpoint.Response, err error) {
	// verify plugin is ready
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		fabrics := p.engine.Fabrics()
		var value []byte

		if fabrics == nil {
			v, e := p.Serialize([]string{})
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		} else {
			v, e := p.Serialize(fabrics)
			if e != nil {
				return nil, nil, e
			}
			value = []byte(v)
		}

		return headers, value, nil
	}, "maas plugin get_fabrics response")
}

// GetFabric returns a fabric from the MAAS server
func (p *MAASPlugin) GetFabric(in *endpoint.Request) (out *endpoint.Response, err error) {
	if err = p.verifyReady(); err != nil {
		return nil, err
	}

	return utility.PluginHandleWrapperWithHeaders(in, func() (map[string][]string, []byte, error) {
		headers := map[string][]string{
			"X-PLUGIN-NAME": {p.Name()},
		}

		names := in.Parameters["name"]
		ids := in.Parameters["id"]
		if names == nil && ids == nil {
			return nil, nil, fmt.Errorf("missing required query parameter. Must specify ?name=_fabric_name_ or ?id=_fabric_id_")
		}

		var matches []*structs.Fabric
		if fabrics := p.engine.Fabrics(); fabrics != nil {
			for _, f := range fabrics {
				for _, name := range names {
					if f.Name == name {
						matches = append(matches, f)
					}
				}
				for _, id := range ids {
					if strconv.Itoa(f.ID) == id {
						matches = append(matches, f)
					}
				}
			}
		}

		if len(matches) == 0 {
			return nil, nil, fmt.Errorf("fabric not found")
		}

		v, e := p.Serialize(matches)
		if e != nil {
			return nil, nil, e
		}
		value := []byte(v)

		return headers, value, nil
	}, "maas plugin get_fabric response")
}

func (p *MAASPlugin) getMatchingMachines(in *endpoint.Request) (matches []*structs.Machine, err error) {
	hosts := in.Parameters["host"]
	ids := in.Parameters["id"]
	if hosts == nil && ids == nil {
		return nil, fmt.Errorf("missing required query parameter. Must specify ?host=_hostname_ or ?id=_system_id_")
	}

	matches = p.findMachines(ids, hosts)

	if len(matches) == 0 {
		p.Log(plugins.LevelWarning, "no matching machines found", map[string]string{"ids": fmt.Sprintf("%v", ids), "hosts": fmt.Sprintf("%v", hosts)})
		return nil, nil
	}

	return matches, nil
}

func (p *MAASPlugin) findMachines(ids []string, hosts []string) []*structs.Machine {
	results := make([]*structs.Machine, 0)

	if len(ids) > 0 {
		for _, id := range ids {
			machine := p.getMachineByID(id)
			if machine != nil {
				results = append(results, machine)
			}
		}
	}

	if len(hosts) > 0 {
		for _, host := range hosts {
			machine := p.getMachineByHost(host)
			if machine != nil {
				results = append(results, machine)
			}
		}
	}

	return results
}

func (p *MAASPlugin) getMachineByID(id string) *structs.Machine {
	if p.engine != nil {
		for _, machine := range p.engine.Machines() {
			if machine.SystemID == id {
				return machine
			}
		}
	}
	return nil
}

func (p *MAASPlugin) getMachineByHost(host string) *structs.Machine {
	if p.engine != nil {
		for _, machine := range p.engine.Machines() {
			if machine.Hostname == host {
				return machine
			}
		}
	}
	return nil
}

func (p MAASPlugin) verifyReady() error {
	if p.engine == nil || !p.engine.Running() {
		return fmt.Errorf("plugin has not been started")
	}

	return nil
}
