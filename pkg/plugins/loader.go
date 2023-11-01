package plugins

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/BGrewell/go-execute"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/pkg/plugins/utility"
	"io/ioutil"
	"strconv"
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
func NewPluginLoader(pluginDirectory string, pluginRoot string, plugConfigs map[string]*PluginConfig, loadUnconfiguredPlugins bool) PluginLoader {
	l := &DefaultPluginLoader{
		PluginDirectory:         pluginDirectory,
		PluginConfigs:           plugConfigs,
		loadUnconfiguredPlugins: loadUnconfiguredPlugins,
		plugins:                 make(map[string]*PluginInfo),
		routeMap:                make(map[string]*HandlerEntry),
		pluginRoot:              pluginRoot,
	}

	return l
}

// executePlugin is called to launch a plugin
func executePlugin(config *PluginConfig) (info *PluginInfo, err error) {
	// If the hash was configured with a sha1 hash verify that the loaded image matches the expected value
	if config.Hash != "" {
		if err := hashValid(config.PluginPath, config.Hash); err == nil {
			return nil, err
		}
	}

	// Execute the plugin
	// TODO: to support execution under another user go-execute has v2 which currently has Linux working but not
	// Windows at this point. Once Windows is supported this will be updated to use the new version and will support
	// setting the user to execute the plugin as.
	stdout, _, exitChan, cancel, err := execute.ExecuteAsyncWithCancel(config.PluginPath, &[]string{"DTAC_PLUGINS=true"})
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	scanner.Scan()
	output := scanner.Text()
	output = strings.Replace(output, "CONNECT{{", "", 1)
	output = strings.Replace(output, "}}", "", 1)
	fields := strings.Split(output, ":")
	port, err := strconv.Atoi(fields[4])
	if err != nil {
		return nil, err
	}

	key, err := utility.DecodeKeyString(fields[6])
	if err != nil {
		return nil, err
	}

	info = &PluginInfo{
		Path:         config.PluginPath,
		Name:         fields[0],
		RootPath:     fields[1],
		Pid:          0,
		Proto:        fields[2],
		Ip:           fields[3],
		Port:         port,
		ApiVersion:   fields[5],
		Key:          key,
		Rpc:          nil,
		CancelToken:  &cancel,
		ExitChan:     exitChan,
		ExitCode:     0,
		PluginConfig: config,
	}

	// Create a go routine to watch for the plugin to exit
	go func() {
		ec := <-info.ExitChan
		info.ExitCode = ec
		info.HasExited = true
	}()
	return info, nil
}

// hashValid is called to verify that the hash of the plugin matches the expected value
func hashValid(pluginPath string, expectedHash string) (err error) {
	// Ensure expected hash is lowercase
	expectedHash = strings.ToLower(expectedHash)

	// Read the binary file from disk
	binaryData, err := ioutil.ReadFile(pluginPath)
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
