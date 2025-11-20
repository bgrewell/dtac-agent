package modules

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/BGrewell/go-execute"
	api "github.com/bgrewell/dtac-agent/api/grpc/go"
	"github.com/bgrewell/dtac-agent/pkg/modules/utility"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// DefaultModuleLoader is the default implementation of the ModuleLoader interface
type DefaultModuleLoader struct {
	ModuleDirectory         string
	ModuleConfigs           map[string]*ModuleConfig
	loadUnconfiguredModules bool
	modules                 map[string]*ModuleInfo
	moduleRoot              string
	tlsCertFile             *string
	tlsKeyFile              *string
	tlsCAFile               *string
	defaultSecure           bool
	logger                  *zap.Logger
}

// Initialize is used to initialize the module loader
func (ml *DefaultModuleLoader) Initialize(secure bool) (loadedModules []*ModuleInfo, err error) {

	// Set default secure
	ml.defaultSecure = secure

	// List modules
	loadedModules = make([]*ModuleInfo, 0)
	mods, err := ml.ListModules()
	if err != nil {
		return nil, err
	}
	ml.logger.Info("enumerated modules", zap.Strings("modules", mods))
	ml.logger.Debug("module configs", zap.Any("configs", ml.ModuleConfigs))
	for _, mod := range mods {
		// Inside here we don't return errors because we want to continue loading other modules. Instead, we log the
		// error and continue
		if config, exists := ml.ModuleConfigs[mod]; (exists && config.Enabled) || (!exists && ml.loadUnconfiguredModules) {

			// Launch modules
			info, err := ml.LaunchModule(ml.ModuleConfigs[mod])
			if err != nil {
				ml.logger.Error("failed to launch module", zap.String("module", mod), zap.Error(err))
				continue
			}

			// Register modules
			err = ml.RegisterModule(info.Name)
			if err != nil {
				ml.logger.Error("failed to register module", zap.String("module", mod), zap.Error(err))
				continue
			}
			loadedModules = append(loadedModules, info)
		}
	}

	return loadedModules, nil
}

// ListModules returns a list of all modules in the module directory
func (ml *DefaultModuleLoader) ListModules() (modules []string, err error) {
	if runtime.GOOS == "windows" {
		return utility.FindModules(ml.ModuleDirectory, "*.module.exe")
	} else if runtime.GOOS == "darwin" {
		return utility.FindModules(ml.ModuleDirectory, "*.module.app")
	} else {
		return utility.FindModules(ml.ModuleDirectory, "*.module")
	}
}

// LaunchModule launches a module and returns the info on the running module
func (ml *DefaultModuleLoader) LaunchModule(config *ModuleConfig) (info *ModuleInfo, err error) {
	info, err = ml.executeModule(config)
	if err != nil {
		return nil, err
	}
	ml.modules[info.Name] = info
	return info, err
}

// RegisterModule registers the module with the agent
func (ml *DefaultModuleLoader) RegisterModule(moduleName string) (err error) {

	if _, ok := ml.modules[moduleName]; !ok {
		return fmt.Errorf("no module was found with the name: %s", moduleName)
	}

	mod := ml.modules[moduleName]

	// Setup security
	creds := insecure.NewCredentials()
	if ml.tlsCertFile != nil && ml.tlsKeyFile != nil && ml.tlsCAFile != nil {
		// Load the certificates from disk
		cert, err := tls.LoadX509KeyPair(*ml.tlsCertFile, *ml.tlsKeyFile)
		if err != nil {
			return fmt.Errorf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := os.ReadFile(*ml.tlsCAFile)
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
	moduleAddress := fmt.Sprintf("%s:%d", mod.IP, mod.Port)
	conn, err := grpc.Dial(moduleAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}
	mod.RPC = api.NewModuleServiceClient(conn)

	// Set up the logging stream
	stream, err := mod.RPC.LoggingStream(context.Background(), &api.LoggingArgs{})
	if err != nil {
		return fmt.Errorf("error setting up module logging: %v", err)
	}
	modLogger := ml.logger.With(zap.String("module", mod.Name))
	go handleLoggingRequests(stream, modLogger)

	// Set up the configuration input
	configJSON, err := json.Marshal(ml.modules[moduleName].ModuleConfig.Config)
	if err != nil {
		// Handle error
		return err
	}

	ra := &api.ModuleRegisterRequest{
		Config:        string(configJSON),
		DefaultSecure: ml.defaultSecure,
	}

	// Call the module's register function
	reply, err := mod.RPC.Register(context.Background(), ra)
	if err != nil {
		return err
	}

	// Store module metadata
	mod.ModuleType = reply.ModuleType
	mod.Capabilities = reply.Capabilities

	return nil
}

func handleLoggingRequests(stream api.ModuleService_LoggingStreamClient, modLogger *zap.Logger) {
	for {
		logMsg, err := stream.Recv()
		if err == io.EOF {
			modLogger.Warn("reached end of log stream. logging for this module has terminated", zap.String("source", "loader"))
			return
		}
		if err != nil {
			modLogger.Error("failed to receive log message. logging for this module has terminated", zap.Error(err), zap.String("source", "loader"))
			return
		}

		fields := make([]zap.Field, 0)
		for _, field := range logMsg.Fields {
			fields = append(fields, zap.Any(field.Key, field.Value))
		}

		switch logMsg.Level {
		case api.LogLevel_DEBUG:
			fields = append(fields, zap.String("source", "module"))
			modLogger.Debug(logMsg.Message, fields...)
		case api.LogLevel_INFO:
			fields = append(fields, zap.String("source", "module"))
			modLogger.Info(logMsg.Message, fields...)
		case api.LogLevel_WARNING:
			fields = append(fields, zap.String("source", "module"))
			modLogger.Warn(logMsg.Message, fields...)
		case api.LogLevel_ERROR:
			fields = append(fields, zap.String("source", "module"), zap.Bool("fatal", false))
			modLogger.Error(logMsg.Message, fields...)
		case api.LogLevel_FATAL:
			fields = append(fields, zap.String("source", "module"), zap.Bool("fatal", true))
			modLogger.Error(logMsg.Message, fields...)
		default:
			fields = append(fields, zap.String("source", "loader"), zap.String("level", logMsg.Level.String()), zap.String("message", logMsg.Message))
			modLogger.Error("received log message with invalid log level", fields...)
		}
	}
}

// UnregisterModule is used to unregister the module
func (ml *DefaultModuleLoader) UnregisterModule(moduleName string) (err error) {
	if _, ok := ml.modules[moduleName]; !ok {
		return fmt.Errorf("module not found: %s", moduleName)
	}
	return nil
}

// CloseModule is used to stop the module process
func (ml *DefaultModuleLoader) CloseModule(moduleName string) (err error) {
	if _, ok := ml.modules[moduleName]; !ok {
		return fmt.Errorf("module not found: %s", moduleName)
	}

	token := *ml.modules[moduleName].CancelToken
	token()

	return nil
}

func (ml *DefaultModuleLoader) executeModule(config *ModuleConfig) (info *ModuleInfo, err error) {
	// If the hash was configured with a sha1 hash verify that the loaded image matches the expected value
	if config.Hash != "" {
		if err := hashValid(config.ModulePath, config.Hash); err != nil {
			return nil, err
		}
	}

	// ensure that the module is writable only to root or the current process user
	if onlyWriteable, err := utility.IsOnlyWritableByUserOrRoot(config.ModulePath); err != nil {
		return nil, fmt.Errorf("failed to check if module is only writeable by root or self: %v", err)
	} else if !onlyWriteable {
		return nil, fmt.Errorf("module has incorrect file permissions. Only root or the current process user should have write access")
	}

	// Set the environment variables for TLS if configuration is present
	envs := []string{"DTAC_MODULES=true"}

	if ml.tlsCertFile != nil && ml.tlsKeyFile != nil {
		certBytes, err := os.ReadFile(*ml.tlsCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read TLS cert file: %s", err.Error())
		}

		keyBytes, err := os.ReadFile(*ml.tlsKeyFile)
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

	// Execute the module
	ml.logger.Debug("executing module",
		zap.String("module", config.ModulePath),
		zap.Strings("envs", envs))
	stdout, _, exitChan, cancel, err := execute.ExecuteAsyncWithCancel(config.ModulePath, &envs)
	if err != nil {
		return nil, err
	}

	// Wait for command to execute and output to be available
	reader := bufio.NewReader(stdout)
	ready, err := execute.SignalRead(reader, 2*time.Second)
	if err != nil || !ready {
		return nil, fmt.Errorf("failed to read from module: %s", err.Error())
	}

	scanner := bufio.NewScanner(reader)
	scanner.Scan()
	output := scanner.Text()
	ml.logger.Info("module output", zap.String("output", output))
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

	info = &ModuleInfo{
		Path:          config.ModulePath,
		Name:          fields[0],
		RootPath:      fields[1],
		Pid:           0,
		RPCProto:      fields[2],
		Proto:         fields[3],
		IP:            fields[4],
		Port:          port,
		APIVersion:    fields[6],
		ModuleOptions: options,
		RPC:           nil,
		CancelToken:   &cancel,
		ExitChan:      exitChan,
		ExitCode:      0,
		ModuleConfig:  config,
	}

	// if RootPath is empty, set it to the module file name minus ".module"
	if info.RootPath == "" {
		// Remove the directory path
		filename := filepath.Base(config.ModulePath)

		// Remove the .module extension
		info.RootPath = strings.TrimSuffix(filename, ".module")
	}

	// Create a go routine to watch for the module to exit
	go func() {
		ec := <-info.ExitChan
		info.ExitCode = ec
		info.HasExited = true
	}()
	return info, nil
}
