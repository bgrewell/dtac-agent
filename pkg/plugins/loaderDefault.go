package plugins

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BGrewell/go-execute"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/helpers"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	api "github.com/intel-innersource/frameworks.automation.dtac.agent/api/grpc/go"
)

// LoadUnloadArgs is a struct that defines the arguments for the load and unload endpoints
type LoadUnloadArgs struct {
	Name string `json:"name"`
}

// DefaultPluginLoader is the default implementation of the PluginLoader interface
type DefaultPluginLoader struct {
	PluginDirectory         string
	PluginConfigs           map[string]*PluginConfig
	loadUnconfiguredPlugins bool
	plugins                 map[string]*PluginInfo
	routeMap                map[string]*HandlerEntry
	endpoints               []*endpoint.Endpoint
	pluginRoot              string
	tlsCertFile             *string
	tlsKeyFile              *string
	tlsCAFile               *string
	defaultSecure           bool
	logger                  *zap.Logger
}

// Initialize is used to initialize the plugin loader
func (pl *DefaultPluginLoader) Initialize(secure bool) (loadedPlugins []*PluginInfo, err error) {

	// Set default secure
	pl.defaultSecure = secure

	// List plugins
	loadedPlugins = make([]*PluginInfo, 0)
	plugs, err := pl.ListPlugins()
	if err != nil {
		return nil, err
	}

	for _, plug := range plugs {
		// Inside here we don't return errors because we want to continue loading other plugins. Instead, we log the
		// error and continue
		if config, exists := pl.PluginConfigs[plug]; (exists && config.Enabled) || (!exists && pl.loadUnconfiguredPlugins) {

			// Launch plugins
			info, err := pl.LaunchPlugin(pl.PluginConfigs[plug])
			if err != nil {
				pl.logger.Error("failed to launch plugin", zap.String("plugin", plug), zap.Error(err))
				continue
			}

			// Register plugins
			err = pl.RegisterPlugin(info.Name)
			if err != nil {
				pl.logger.Error("failed to register plugin", zap.String("plugin", plug), zap.Error(err))
				continue
			}
			loadedPlugins = append(loadedPlugins, info)
		}
	}

	authz := endpoint.AuthGroupAdmin.String()
	// Register control routes GET methods are just there for ease of use
	endpoints := []*endpoint.Endpoint{
		endpoint.NewEndpoint("load", endpoint.ActionRead, "load a plugin", pl.Load, secure, authz, endpoint.WithParameters(LoadUnloadArgs{})),
		endpoint.NewEndpoint("load", endpoint.ActionCreate, "load a plugin", pl.Load, secure, authz, endpoint.WithParameters(LoadUnloadArgs{})),
		endpoint.NewEndpoint("unload", endpoint.ActionRead, "unload a plugin", pl.Unload, secure, authz, endpoint.WithParameters(LoadUnloadArgs{})),
		endpoint.NewEndpoint("unload", endpoint.ActionCreate, "unload a plugin", pl.Unload, secure, authz, endpoint.WithParameters(LoadUnloadArgs{})),
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
	info, err = pl.executePlugin(config)
	if err != nil {
		return nil, err
	}
	pl.plugins[info.Name] = info
	return info, err
}

// RegisterPlugin registers the plugin routes with Gin
func (pl *DefaultPluginLoader) RegisterPlugin(pluginName string) (err error) {

	if _, ok := pl.plugins[pluginName]; !ok {
		return fmt.Errorf("no plugin was found with the name: %s", pluginName)
	}

	plug := pl.plugins[pluginName]

	// Setup security
	creds := insecure.NewCredentials()
	if pl.tlsCertFile != nil && pl.tlsKeyFile != nil && pl.tlsCAFile != nil {
		// Load the certificates from disk
		cert, err := tls.LoadX509KeyPair(*pl.tlsCertFile, *pl.tlsKeyFile)
		if err != nil {
			return fmt.Errorf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := os.ReadFile(*pl.tlsCAFile)
		if err != nil {
			return fmt.Errorf("could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return fmt.Errorf("failed to append ca certs")
		}

		creds = credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      certPool,
		})
	}

	// Connect the rpc client
	pluginAddress := fmt.Sprintf("%s:%d", plug.IP, plug.Port)
	conn, err := grpc.Dial(pluginAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	plug.RPC = api.NewPluginServiceClient(conn)

	// Set up the logging stream
	stream, err := plug.RPC.LoggingStream(context.Background(), &api.LoggingArgs{})
	if err != nil {
		return fmt.Errorf("error setting up plugin logging: %v", err)
	}
	plugLogger := pl.logger.With(zap.String("plugin", plug.Name))
	go handleLoggingRequests(stream, plugLogger)

	// Set up the configuration input
	configJSON, err := json.Marshal(pl.plugins[pluginName].PluginConfig.Config)
	if err != nil {
		// Handle error
		return err
	}

	ra := &api.RegisterRequest{
		Config:        string(configJSON),
		DefaultSecure: pl.defaultSecure,
	}

	// Call the plugins register function
	reply, err := plug.RPC.Register(context.Background(), ra)
	if err != nil {
		return err
	}

	// Convert the endpoints
	plug.Endpoints = reply.Endpoints

	// Record routes
	endpoints := make([]*endpoint.Endpoint, 0)
	for _, aep := range plug.Endpoints {
		// Ensure action is valid
		_, err := endpoint.ParseAction(aep.Action)
		if err != nil {
			return err
		}

		// Register endpoints
		handleFuncName := aep.Path
		aep.Path = path.Join(plug.RootPath, aep.Path)
		ep := utility.ConvertPluginEndpointToEndpoint(aep)
		endpoints = append(endpoints, ep)

		// Record route map
		key := fmt.Sprintf("%s:%s", aep.Action, aep.Path)
		entry := &HandlerEntry{
			PluginName: plug.Name,
			HandleFunc: handleFuncName,
		}
		pl.routeMap[key] = entry
	}
	pl.endpoints = append(pl.endpoints, endpoints...)

	return nil
}

func handleLoggingRequests(stream api.PluginService_LoggingStreamClient, plugLogger *zap.Logger) {
	for {
		logMsg, err := stream.Recv()
		if err == io.EOF {
			plugLogger.Warn("reached end of log stream. logging for this plugin has terminated", zap.String("source", "loader"))
			return
		}
		if err != nil {
			plugLogger.Error("failed to receive log message. logging for this plugin has terminated", zap.Error(err), zap.String("source", "loader"))
			return
		}

		fields := make([]zap.Field, 0)
		for _, field := range logMsg.Fields {
			fields = append(fields, zap.Any(field.Key, field.Value))
		}

		switch logMsg.Level {
		case api.LogLevel_DEBUG:
			fields = append(fields, zap.String("source", "plugin"))
			plugLogger.Debug(logMsg.Message, fields...)
		case api.LogLevel_INFO:
			fields = append(fields, zap.String("source", "plugin"))
			plugLogger.Info(logMsg.Message, fields...)
		case api.LogLevel_WARNING:
			fields = append(fields, zap.String("source", "plugin"))
			plugLogger.Warn(logMsg.Message, fields...)
		case api.LogLevel_ERROR:
			fields = append(fields, zap.String("source", "plugin"), zap.Bool("fatal", false))
			plugLogger.Error(logMsg.Message, fields...)
		case api.LogLevel_FATAL:
			fields = append(fields, zap.String("source", "plugin"), zap.Bool("fatal", true))
			plugLogger.Error(logMsg.Message, fields...)
		default:
			fields = append(fields, zap.String("source", "loader"), zap.String("level", logMsg.Level.String()), zap.String("message", logMsg.Message))
			plugLogger.Error("received log message with invalid log level", fields...)
		}
	}
}

// UnregisterPlugin is used to unregister the plugin from Gin
func (pl *DefaultPluginLoader) UnregisterPlugin(pluginName string) (err error) {
	if _, ok := pl.plugins[pluginName]; !ok {
		return errors.New("plugin not found")
	}
	return nil
}

// ClosePlugin is used to stop the plugin process
func (pl *DefaultPluginLoader) ClosePlugin(pluginName string) (err error) {
	if _, ok := pl.plugins[pluginName]; !ok {
		return errors.New("plugin not found")
	}

	token := *pl.plugins[pluginName].CancelToken
	token()

	return nil
}

// Load is used to load a plugin by name
func (pl *DefaultPluginLoader) Load(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		if m := in.Parameters["name"]; m[0] != "" {
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

			return []byte("plugin loaded"), nil
		}

		return nil, errors.New("missing 'name' parameter specifying the plugin name")

	}, "plugin loaded")
}

// Unload is used to unload a plugin by name
func (pl *DefaultPluginLoader) Unload(in *endpoint.Request) (out *endpoint.Response, err error) {
	return helpers.HandleWrapper(in, func() ([]byte, error) {
		if m := in.Parameters["name"]; m[0] != "" {
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

			return []byte("plugin unloaded"), nil
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
func (pl *DefaultPluginLoader) CallShim(ep *endpoint.Endpoint, in *endpoint.Request) (out *endpoint.Response, err error) {

	// Extract the RouteKey
	keyPath := strings.TrimLeft(strings.Replace(ep.Path, pl.pluginRoot, "", 1), "/")
	routeKey := fmt.Sprintf("%s:%s", ep.Action, keyPath)

	// Get the HandlerEntry
	if _, ok := pl.routeMap[routeKey]; !ok {
		return nil, errors.New("a handler was not found for the requested resource")
	}

	handler := pl.routeMap[routeKey]

	// Make sure plugin isn't canceled
	if pl.plugins[handler.PluginName].HasExited {
		return nil, fmt.Errorf("the plugin has exited with code: %d", pl.plugins[handler.PluginName].ExitCode)
	}

	// Get the plugin
	plug := pl.plugins[handler.PluginName]

	// Setup the request message
	erm := &api.EndpointRequestMessage{
		Method:  fmt.Sprintf("%s:%s", ep.Action, handler.HandleFunc),
		Request: utility.EndpointRequestToAPIEndpointRequest(in),
	}

	// Make the rpc call
	ret, err := plug.RPC.Call(context.Background(), erm)
	if err != nil {
		return nil, fmt.Errorf("failed to call plugin function: %s", err)
	}
	out = utility.APIEndpointResponseToEndpointResponse(ret.Response)

	return out, nil
}

// executePlugin is called to launch a plugin
func (pl *DefaultPluginLoader) executePlugin(config *PluginConfig) (info *PluginInfo, err error) {
	// If the hash was configured with a sha1 hash verify that the loaded image matches the expected value
	if config.Hash != "" {
		if err := hashValid(config.PluginPath, config.Hash); err == nil {
			return nil, err
		}
	}

	// ensure that the plugin is writable only to root or the current process user
	if onlyWriteable, err := utility.IsOnlyWritableByUserOrRoot(config.PluginPath); err != nil {
		return nil, fmt.Errorf("failed to check if plugin is only writeable by root or self: %v", err)
	} else if !onlyWriteable {
		return nil, fmt.Errorf("plugin has incorrect file permissions. Only root or the current process user should have write access")
	}

	// Set the environment variables for TLS if configuration is present
	envs := []string{"DTAC_PLUGINS=true"}

	if pl.tlsCertFile != nil && pl.tlsKeyFile != nil {
		certBytes, err := os.ReadFile(*pl.tlsCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read TLS cert file: %s", err.Error())
		}

		keyBytes, err := os.ReadFile(*pl.tlsKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read TLS key file: %s", err.Error())
		}
		cert := string(certBytes)
		key := string(keyBytes)

		envs = append(envs, fmt.Sprintf("DTAC_TLS_CERT=%s", cert), fmt.Sprintf("DTAC_TLS_KEY=%s", key))

		certBytes = nil
		keyBytes = nil
		cert = ""
		key = ""
	}

	// Execute the plugin
	// TODO: to support execution under another user go-execute has v2 which currently has Linux working but not
	// Windows at this point. Once Windows is supported this will be updated to use the new version and will support
	// setting the user to execute the plugin as.
	stdout, _, exitChan, cancel, err := execute.ExecuteAsyncWithCancel(config.PluginPath, &envs)
	if err != nil {
		return nil, err
	}

	// Wait for command to execute and output to be available - TODO: Expose timeout as a config option
	reader := bufio.NewReader(stdout)
	ready, err := execute.SignalRead(reader, 2*time.Second)
	if err != nil || !ready {
		return nil, fmt.Errorf("failed to read from plugin: %s", err.Error())
	}

	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	output := scanner.Text()
	pl.logger.Info("plugin output", zap.String("output", output))
	output = strings.Replace(output, "CONNECT{{", "", 1)
	output = strings.Replace(output, "}}", "", 1)
	fields := strings.Split(output, ":")
	port, err := strconv.Atoi(fields[5])
	if err != nil {
		return nil, err
	}

	options, err := ParseOptions(fields[7])
	if err != nil {
		return nil, err
	}

	// Try to clean up any key material that may be in memory
	envs = nil
	runtime.GC()

	info = &PluginInfo{
		Path:          config.PluginPath,
		Name:          fields[0],
		RootPath:      fields[1],
		Pid:           0,
		RPCProto:      fields[2],
		Proto:         fields[3],
		IP:            fields[4],
		Port:          port,
		APIVersion:    fields[6],
		PluginOptions: options,
		RPC:           nil,
		CancelToken:   &cancel,
		ExitChan:      exitChan,
		ExitCode:      0,
		PluginConfig:  config,
	}

	// if RootPath is empty, set it to the plugin file name minus ".plugin"
	if info.RootPath == "" {
		// Remove the directory path
		filename := filepath.Base(config.PluginPath)

		// Remove the .plugin extension
		info.RootPath = strings.TrimSuffix(filename, ".plugin")
	}

	// Create a go routine to watch for the plugin to exit
	go func() {
		ec := <-info.ExitChan
		info.ExitCode = ec
		info.HasExited = true
	}()
	return info, nil
}
