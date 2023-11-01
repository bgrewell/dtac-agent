package plugins

import (
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"io"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

// PluginHost is the default interface for the plugin host
type DefaultPluginHost struct {
	Plugin     Plugin
	Proto      string
	Ip         string
	ApiVersion string
	port       int
	encryptor  *utility.RpcEncryptor
}

// Serve starts the plugin host
func (ph *DefaultPluginHost) Serve() error {
	// Hacky way to keep the net.rpc package from complaining about some method signatures
	logger := log.Default()
	logger.SetOutput(io.Discard)

	// Verify that the ENV variable is set else exit with helpful message
	if os.Getenv("DTAC_PLUGINS") == "" {
		fmt.Println("============================ WARNING ============================")
		fmt.Println("This is a DTAC plugin and is not designed to be executed directly")
		fmt.Println("Please use the DTAC agent to load this plugin")
		fmt.Println("==================================================================")
	}

	// Register plugin
	err := rpc.Register(ph.Plugin)
	logger.SetOutput(os.Stderr)
	if err != nil {
		return err
	}

	// Find a TCP port to use
	ph.port, err = utility.GetUnusedTcpPort()
	if err != nil {
		return err
	}

	// Output connection information ( format: CONNECT{{NAME:ROUTE_ROOT:PROTO:IP:PORT:VER:SYM_KEY}} )
	fmt.Printf("CONNECT{{%s:%s:%s:%s:%d:%s:%s}}\n", ph.Plugin.Name(), ph.Plugin.RootPath(), ph.Proto, ph.Ip, ph.port, ph.ApiVersion, ph.encryptor.KeyString())

	// Listen for connections
	l, e := net.Listen(ph.Proto, fmt.Sprintf("%s:%d", ph.Ip, ph.port))
	if e != nil {
		return e
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go jsonrpc.ServeConn(conn) // Serve connection with JSON-RPC
	}
}

// GetPort returns the port the plugin host is listening on
func (ph *DefaultPluginHost) GetPort() int {
	return ph.port
}
