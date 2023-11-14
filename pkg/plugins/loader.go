package plugins

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/types/endpoint"
	"go.uber.org/zap"
	"os"
	"strings"
)

// PluginLoader is the interface that any plugin loaders must satisfy. The default concrete implementation is
// the DefaultPluginLoader which at this point in time is the only planned loader.
type PluginLoader interface {
	Initialize(secure bool) (loadedPlugins []*PluginInfo, err error)
	ListPlugins() (plugins []string, err error)
	LaunchPlugin(config *PluginConfig) (info *PluginInfo, err error)
	RegisterPlugin(pluginName string) (err error)
	UnregisterPlugin(pluginName string) (err error)
	ClosePlugin(pluginName string) (err error)
	Endpoints() []*endpoint.Endpoint
	CallShim(ep *endpoint.Endpoint, in *endpoint.InputArgs) (out *endpoint.ReturnVal, err error)
}

// NewPluginLoader takes in the plugin directory, the sanity cookie and the routeGroup.
//
//	pluginDirectory: the directory that contains all the plugins
//	cookie: the sanity cookie that is used to verify that what is being executed is the expected plugin
//	routeGroup: the routeGroup that all the plugin routes will be placed inside
func NewPluginLoader(pluginDirectory string, pluginRoot string, plugConfigs map[string]*PluginConfig, loadUnconfiguredPlugins bool, tlsCertFile *string, tlsKeyFile *string, tlsCAFile *string, logger *zap.Logger) PluginLoader {
	l := &DefaultPluginLoader{
		PluginDirectory:         pluginDirectory,
		PluginConfigs:           plugConfigs,
		loadUnconfiguredPlugins: loadUnconfiguredPlugins,
		plugins:                 make(map[string]*PluginInfo),
		routeMap:                make(map[string]*HandlerEntry),
		pluginRoot:              pluginRoot,
		tlsCertFile:             tlsCertFile,
		tlsKeyFile:              tlsKeyFile,
		tlsCAFile:               tlsCAFile,
		logger:                  logger,
	}

	return l
}

// hashValid is called to verify that the hash of the plugin matches the expected value
func hashValid(pluginPath string, expectedHash string) (err error) {
	// Ensure expected hash is lowercase
	expectedHash = strings.ToLower(expectedHash)

	// Read the binary file from disk
	binaryData, err := os.ReadFile(pluginPath)
	if err != nil {
		return
	}

	// Compute the SHA1 hash of the binary data
	sha1Hash := sha1.Sum(binaryData)

	// Convert the hash to a hex-encoded string
	sha1Hex := hex.EncodeToString(sha1Hash[:])

	if sha1Hex != expectedHash {
		return fmt.Errorf("SHA1 hash: %s did not match expected value %s\"", sha1Hex, expectedHash)
	}

	return nil
}
