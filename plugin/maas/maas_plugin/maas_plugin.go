package maas_plugin

import (
	"fmt"
	"github.com/BGrewell/dtac-agent/plugin/maas/maas_plugin/engine"
	"github.com/BGrewell/dtac-agent/plugin/maas/maas_plugin/structs"
	plugins "github.com/bgrewell/gin-plugins"
	"net/http"
	"reflect"
)

// Ensure that our type meets the requirements for being a plugin
var _ plugins.Plugin = &MAASPlugin{}

type MAASPlugin struct {
	// PluginBase provides some helper functions
	plugins.PluginBase
	settings *structs.MAASSettings
	engine   *engine.Engine
}

func (p MAASPlugin) RouteRoot() string {
	return "maas"
}

func (p MAASPlugin) Name() string {
	t := reflect.TypeOf(p)
	return t.Name()
}

func (p *MAASPlugin) Register(args plugins.RegisterArgs, reply *plugins.RegisterReply) error {
	routes := make([]*plugins.Route, 0)

	// Add the start function
	routes = append(routes, &plugins.Route{
		Path:       "start",
		Method:     http.MethodGet,
		HandleFunc: "Start",
	})

	// Add the stop function
	routes = append(routes, &plugins.Route{
		Path:       "stop",
		Method:     http.MethodGet,
		HandleFunc: "Stop",
	})

	// Add the status function
	routes = append(routes, &plugins.Route{
		Path:       "status",
		Method:     http.MethodGet,
		HandleFunc: "Status",
	})

	// Add the machines function
	routes = append(routes, &plugins.Route{
		Path:       "machines",
		Method:     http.MethodGet,
		HandleFunc: "GetMachines",
	})

	// Add the ids function
	routes = append(routes, &plugins.Route{
		Path:       "machines/ids",
		Method:     http.MethodGet,
		HandleFunc: "GetMachinesIds",
	})

	// Add the machine function
	routes = append(routes, &plugins.Route{
		Path:       "machine",
		Method:     http.MethodGet,
		HandleFunc: "GetMachine",
	})

	*reply = plugins.RegisterReply{Routes: make([]*plugins.Route, 1)}
	reply.Routes = routes

	// Return no error
	return nil
}

func (p *MAASPlugin) Start(args plugins.Args, c *string) error {
	// Pass creds and start the engine to poll for information
	p.settings = &structs.MAASSettings{
		Server:          "maas.edge.lan",
		ConsumerToken:   "bLpahCkJqVtyjymhjq",
		AuthToken:       "fRZfhr8ngu2DbwAvdD",
		AuthSignature:   "n4nda7VhrAngL9ZjxYL2NaGZMAQFtuUH",
		MachinePollSecs: 30,
	} // TODO: This should all be in the configuration file plugin section

	p.engine = &engine.Engine{
		Settings: p.settings,
	}

	err := p.engine.Start()
	if err != nil {
		return err
	}

	return fmt.Errorf("started ... but this function is hardcoded for testing and needs to be finished")
}

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
		status := structs.Status{
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

func (p *MAASPlugin) GetMachines(args plugins.Args, c *string) error {
	if p.engine == nil {
		return fmt.Errorf("plugin has not been started")
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

func (p *MAASPlugin) GetMachinesIds(args plugins.Args, c *string) error {
	if p.engine == nil {
		return fmt.Errorf("plugin has not been started")
	}

	type MachinePlug struct {
		Hostname string `json:"hostname"`
		SystemId string `json:"system_id"`
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
				SystemId: machine.SystemId,
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

func (p *MAASPlugin) GetMachine(args plugins.Args, c *string) error {
	hosts := args.QueryParams["host"]
	ids := args.QueryParams["id"]
	if len(hosts) == 0 && len(ids) == 0 {
		fmt.Errorf("missing required query parameter. Must specify ?host=<hostname> or ?id=<system_id>")
	}
	return fmt.Errorf("not implemented: host = %s id = %s", args.QueryParams["host"], args.QueryParams["id"])
}

func (p *MAASPlugin) GetMachineId(args plugins.Args, c *string) error {
	hosts := args.QueryParams["host"]
	if len(hosts) == 0 {
		fmt.Errorf("missing required query parameter ?host=<hostname>")
	}
	return fmt.Errorf("not implemented: machine = %s", args.QueryParams["host"])
}

func (p *MAASPlugin) GetMachinePool(args plugins.Args, c *string) error {
	hosts := args.QueryParams["host"]
	ids := args.QueryParams["id"]
	if len(hosts) == 0 && len(ids) == 0 {
		fmt.Errorf("missing required query parameter. Must specify ?host=<hostname> or ?id=<system_id>")
	}
	return fmt.Errorf("not implemented: host = %s id = %s", args.QueryParams["host"], args.QueryParams["id"])
}

func (p *MAASPlugin) GetMachineUnidentitiedInterfaces(args plugins.Args, c *string) error {
	hosts := args.QueryParams["host"]
	ids := args.QueryParams["id"]
	if len(hosts) == 0 && len(ids) == 0 {
		fmt.Errorf("missing required query parameter. Must specify ?host=<hostname> or ?id=<system_id>")
	}
	return fmt.Errorf("not implemented: host = %s id = %s", args.QueryParams["host"], args.QueryParams["id"])
}

func (p *MAASPlugin) GetFabrics(args plugins.Args, c *string) error {
	return fmt.Errorf("not implemented: host = %s id = %s", args.QueryParams["host"], args.QueryParams["id"])
}

func (p *MAASPlugin) GetFabric(args plugins.Args, c *string) error {
	return fmt.Errorf("not implemented: host = %s id = %s", args.QueryParams["host"], args.QueryParams["id"])
}

//func (p *MAASPlugin) Hello(args plugins.Args, c *string) error {
//	v, e := p.Serialize(p.message)
//	if e != nil {
//		return e
//	}
//	*c = v
//	return nil
//}
