package maasplugin

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maasplugin/engine"
	structs2 "github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/maas/maasplugin/structs"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	plugins "github.com/bgrewell/gin-plugins"
)

// Ensure that our type meets the requirements for being a plugin
var _ plugins.Plugin = &MAASPlugin{}

// MAASPlugin is the main plugin type
type MAASPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	settings *structs2.MAASSettings
	engine   *engine.Engine
}

// RouteRoot returns the root path for the plugin
func (p MAASPlugin) RouteRoot() string {
	return "maas"
}

// Name returns the name of the plugin
func (p MAASPlugin) Name() string {
	t := reflect.TypeOf(p)
	return t.Name()
}

// Register registers the plugin with the plugin manager
func (p *MAASPlugin) Register(args plugins.RegisterArgs, reply *plugins.RegisterReply) error {
	routes := make([]*plugins.Route, 0)

	// Open the log file for writing
	if val, ok := args.Config["logfile"]; ok {
		file, err := os.OpenFile(val.(string), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Failed to open log file:", err)
		}

		// Set the log output to the file
		log.SetOutput(file)

	}

	server := "localhost"
	if _, ok := args.Config["server"]; ok {
		server = args.Config["server"].(string)
	}
	pollInt := 30
	if _, ok := args.Config["poll_interval_secs"]; ok {
		pollInt = args.Config["poll_interval_secs"].(int)
	}
	required := []string{"consumer_token", "auth_token", "auth_signature"}
	for _, requiredKey := range required {
		if _, ok := args.Config[requiredKey]; !ok {
			return fmt.Errorf("missing requried plugin configuration: %s", requiredKey)
		}
	}

	// Setup credentials
	p.settings = &structs2.MAASSettings{
		Server:          server,
		ConsumerToken:   args.Config["consumer_token"].(string),
		AuthToken:       args.Config["auth_token"].(string),
		AuthSignature:   args.Config["auth_signature"].(string),
		MachinePollSecs: pollInt,
	}

	// Define the routes
	routeDefs := []struct {
		path       string
		method     string
		handleFunc string
	}{
		{"start", http.MethodGet, "Start"},
		{"stop", http.MethodGet, "Stop"},
		{"status", http.MethodGet, "Status"},
		{"machines", http.MethodGet, "GetMachines"},
		{"machines/ids", http.MethodGet, "GetMachinesIDs"},
		{"machines/pools", http.MethodGet, "GetMachinesPools"},
		{"machines/status", http.MethodGet, "GetMachinesStatuses"},
		{"machines/interfaces", http.MethodGet, "GetMachinesInterfaces"},
		{"machine", http.MethodGet, "GetMachine"},
		{"machine/id", http.MethodGet, "GetMachineID"},
		{"machine/pool", http.MethodGet, "GetMachinePool"},
		{"machine/status", http.MethodGet, "GetMachineStatus"},
		{"machine/interfaces", http.MethodGet, "GetMachineInterfaces"},
		{"fabrics", http.MethodGet, "GetFabrics"},
		{"fabric", http.MethodGet, "GetFabric"},
	}

	// Create the routes
	for _, def := range routeDefs {
		routes = append(routes, newRoute(def.path, def.method, def.handleFunc))
	}

	*reply = plugins.RegisterReply{Routes: make([]*plugins.Route, 1)}
	reply.Routes = routes

	var startReply string
	if _, ok := args.Config["auto_start"]; ok {
		return p.Start(plugins.Args{}, &startReply)
	}

	// Return no error
	return nil
}

// Start starts the plugin
func (p *MAASPlugin) Start(args plugins.Args, c *string) error {

	if p.engine == nil || !p.engine.Running() {
		p.engine = &engine.Engine{
			Settings: p.settings,
		}
	} else {
		return fmt.Errorf("engine is already running. you must stop first")
	}

	err := p.engine.Start()
	if err != nil {
		return err
	}

	*c = "{\"status\": \"started\"}"
	return nil
}

// Stop stops the plugin
func (p *MAASPlugin) Stop(args plugins.Args, c *string) error {
	// Clear creds and stop the engine
	if p.engine.Running() {
		err := p.engine.Stop()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("plugin engine is not running")
	}
	return fmt.Errorf("stopped ... but this function is hardcoded for testing and needs to be finished")
}

// Status returns the status of the plugin
func (p *MAASPlugin) Status(args plugins.Args, c *string) error {
	if p.engine != nil {
		err := ""
		if p.engine.ErrDetails() != nil {
			err = p.engine.ErrDetails().Error()
		}

		count := 0
		machines := p.engine.Machines()
		if machines != nil {
			count = len(machines)
		}
		status := structs2.Status{
			Running:      p.engine.Running(),
			Failed:       p.engine.Failed(),
			ErrDetails:   err,
			MachineCount: count,
		}
		v, e := p.Serialize(status)
		if e != nil {
			return e
		}
		*c = v
	} else {
		v, e := p.Serialize("plugin has not been started")
		if e != nil {
			return e
		}
		*c = v
	}
	return nil
}

// GetMachines returns a list of machines from the MAAS server
func (p *MAASPlugin) GetMachines(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	machines := p.engine.Machines()
	if machines == nil {
		v, e := p.Serialize([]string{})
		if e != nil {
			return e
		}
		*c = v
	} else {
		v, e := p.Serialize(machines)
		if e != nil {
			return e
		}
		*c = v
	}
	return nil
}

// GetMachinesIDs returns a list of machine ids from the MAAS server
func (p *MAASPlugin) GetMachinesIDs(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	type MachinePlug struct {
		Hostname string `json:"hostname"`
		SystemID string `json:"system_id"`
	}

	machines := p.engine.Machines()
	if machines == nil {
		v, e := p.Serialize([]MachinePlug{})
		if e != nil {
			return e
		}
		*c = v
	} else {
		ids := make([]MachinePlug, 0)
		for _, machine := range machines {
			ids = append(ids, MachinePlug{
				Hostname: machine.Hostname,
				SystemID: machine.SystemID,
			})
		}
		v, e := p.Serialize(ids)
		if e != nil {
			return e
		}
		*c = v
	}
	return nil
}

// GetMachinesPools returns a list of machine pools from the MAAS server
func (p *MAASPlugin) GetMachinesPools(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	type Result struct {
		Machine string      `json:"machine"`
		Pool    interface{} `json:"pool"`
	}
	results := make([]Result, 0)

	for _, match := range p.engine.Machines() {
		r := Result{
			Machine: match.Hostname,
			Pool:    match.Pool,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

// GetMachinesStatuses returns a list of machine statuses from the MAAS server
func (p *MAASPlugin) GetMachinesStatuses(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	type Result struct {
		Machine string      `json:"machine"`
		Status  interface{} `json:"status"`
	}
	results := make([]Result, 0)

	for _, match := range p.engine.Machines() {
		r := Result{
			Machine: match.Hostname,
			Status:  match.Status,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

// GetMachinesInterfaces returns a list of machine interfaces from the MAAS server
func (p *MAASPlugin) GetMachinesInterfaces(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	type Result struct {
		Machine    string                     `json:"machine"`
		Interfaces []structs2.InterfaceStruct `json:"interfaces"`
	}
	results := make([]Result, 0)

	for _, match := range p.engine.Machines() {
		r := Result{
			Machine:    match.Hostname,
			Interfaces: match.InterfaceSet,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

// GetMachine returns a machine from the MAAS server
func (p *MAASPlugin) GetMachine(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	matches, err := p.getMatchingMachines(args)
	if err != nil {
		return err
	}

	type Result struct {
		Machine *structs2.Machine `json:"machine"`
	}
	results := make([]Result, 0)

	for _, match := range matches {
		r := Result{
			Machine: match,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

// GetMachineID returns a machine id from the MAAS server
func (p *MAASPlugin) GetMachineID(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	matches, err := p.getMatchingMachines(args)
	if err != nil {
		return err
	}

	type Result struct {
		Machine string `json:"machine"`
		ID      string `json:"id"`
	}
	results := make([]Result, 0)

	for _, match := range matches {
		r := Result{
			Machine: match.Hostname,
			ID:      match.SystemID,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

// GetMachinePool returns a machine pool from the MAAS server
func (p *MAASPlugin) GetMachinePool(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	matches, err := p.getMatchingMachines(args)
	if err != nil {
		return err
	}

	type Result struct {
		Machine string      `json:"machine"`
		Pool    interface{} `json:"pool"`
	}
	results := make([]Result, 0)

	for _, match := range matches {
		r := Result{
			Machine: match.Hostname,
			Pool:    match.Pool,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

// GetMachineStatus returns a machine status from the MAAS server
func (p *MAASPlugin) GetMachineStatus(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	matches, err := p.getMatchingMachines(args)
	if err != nil {
		return err
	}

	type Result struct {
		Machine string      `json:"machine"`
		Status  interface{} `json:"status"`
	}
	results := make([]Result, 0)

	for _, match := range matches {
		r := Result{
			Machine: match.Hostname,
			Status:  match.Status,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

// GetMachineInterfaces returns a machine interfaces from the MAAS server
func (p *MAASPlugin) GetMachineInterfaces(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	matches, err := p.getMatchingMachines(args)
	if err != nil {
		return err
	}

	type Result struct {
		Machine    string                     `json:"machine"`
		Interfaces []structs2.InterfaceStruct `json:"interfaces"`
	}
	results := make([]Result, 0)

	for _, match := range matches {
		r := Result{
			Machine:    match.Hostname,
			Interfaces: match.InterfaceSet,
		}
		results = append(results, r)
	}

	var v string
	var e error
	if len(results) == 1 {
		v, e = p.Serialize(results[0])
	} else {
		v, e = p.Serialize(results)
	}

	if e != nil {
		return e
	}
	*c = v

	return nil
}

func (p *MAASPlugin) getMatchingMachines(args plugins.Args) (matches []*structs2.Machine, err error) {
	hosts := args.QueryParams["host"]
	ids := args.QueryParams["id"]
	if hosts == nil && ids == nil {
		return nil, fmt.Errorf("missing required query parameter. Must specify ?host=_hostname_ or ?id=_system_id_")
	}

	matches = p.findMachines(ids, hosts)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no matching machines found")
	}

	return matches, nil
}

func (p *MAASPlugin) findMachines(ids []string, hosts []string) []*structs2.Machine {
	results := make([]*structs2.Machine, 0)

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

// GetFabrics returns a list of fabrics from the MAAS server
func (p *MAASPlugin) GetFabrics(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	fabrics := p.engine.Fabrics()
	if fabrics == nil {
		v, e := p.Serialize([]string{})
		if e != nil {
			return e
		}
		*c = v
	} else {
		v, e := p.Serialize(fabrics)
		if e != nil {
			return e
		}
		*c = v
	}
	return nil
}

// GetFabric returns a fabric from the MAAS server
func (p *MAASPlugin) GetFabric(args plugins.Args, c *string) error {
	if err := p.verifyReady(args, nil, nil, nil); err != nil {
		return err
	}

	names := args.QueryParams["name"]
	ids := args.QueryParams["id"]
	if names == nil && ids == nil {
		return fmt.Errorf("missing required query parameter. Must specify ?name=_fabric_name_ or ?id=_fabric_id_")
	}

	fabrics := p.engine.Fabrics()
	if fabrics != nil {
		matches := make([]*structs2.Fabric, 0)
		for _, fabric := range fabrics {
			for _, name := range names {
				if fabric.Name == name {
					matches = append(matches, fabric)
				}
			}
			for _, id := range ids {
				if strconv.Itoa(fabric.ID) == id {
					matches = append(matches, fabric)
				}
			}
		}

		if len(matches) > 0 {
			v, e := p.Serialize(matches)
			if e != nil {
				return e
			}
			*c = v
			return nil
		}
	}

	return fmt.Errorf("fabric not found")
}

func (p MAASPlugin) verifyReady(args plugins.Args, queryReqs, headerReqs, bodyReqs *[]string) error {
	if p.engine == nil || !p.engine.Running() {
		return fmt.Errorf("plugin has not been started")
	}

	return nil
}

func (p *MAASPlugin) getMachineByID(id string) *structs2.Machine {
	if p.engine != nil {
		for _, machine := range p.engine.Machines() {
			if machine.SystemID == id {
				return machine
			}
		}
	}
	return nil
}

func (p *MAASPlugin) getMachineByHost(host string) *structs2.Machine {
	if p.engine != nil {
		for _, machine := range p.engine.Machines() {
			if machine.Hostname == host {
				return machine
			}
		}
	}
	return nil
}

func newRoute(path, method, handleFunc string) *plugins.Route {
	return &plugins.Route{
		Path:       path,
		Method:     method,
		HandleFunc: handleFunc,
	}
}
