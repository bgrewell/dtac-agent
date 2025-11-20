package modules

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"go.uber.org/zap"
	"os"
	"strings"
)

// ModuleLoader is the interface that any module loaders must satisfy. The default concrete implementation is
// the DefaultModuleLoader which at this point in time is the only planned loader.
type ModuleLoader interface {
	Initialize(secure bool) (loadedModules []*ModuleInfo, err error)
	ListModules() (modules []string, err error)
	LaunchModule(config *ModuleConfig) (info *ModuleInfo, err error)
	RegisterModule(moduleName string) (err error)
	UnregisterModule(moduleName string) (err error)
	CloseModule(moduleName string) (err error)
	Endpoints() []*endpoint.Endpoint
	CallShim(ep *endpoint.Endpoint, in *endpoint.Request) (out *endpoint.Response, err error)
}

// NewModuleLoader takes in the module directory, the sanity cookie and the routeGroup.
//
//	moduleDirectory: the directory that contains all the modules
//	cookie: the sanity cookie that is used to verify that what is being executed is the expected module
//	routeGroup: the routeGroup that all the module routes will be placed inside
func NewModuleLoader(moduleDirectory string, moduleRoot string, modConfigs map[string]*ModuleConfig, loadUnconfiguredModules bool, tlsCertFile *string, tlsKeyFile *string, tlsCAFile *string, logger *zap.Logger) ModuleLoader {
	l := &DefaultModuleLoader{
		ModuleDirectory:         moduleDirectory,
		ModuleConfigs:           modConfigs,
		loadUnconfiguredModules: loadUnconfiguredModules,
		modules:                 make(map[string]*ModuleInfo),
		moduleRoot:              moduleRoot,
		tlsCertFile:             tlsCertFile,
		tlsKeyFile:              tlsKeyFile,
		tlsCAFile:               tlsCAFile,
		logger:                  logger,
	}

	return l
}

// hashValid is called to verify that the hash of the module matches the expected value
func hashValid(modulePath string, expectedHash string) (err error) {
	// Ensure expected hash is lowercase
	expectedHash = strings.ToLower(expectedHash)

	// Read the binary file from disk
	binaryData, err := os.ReadFile(modulePath)
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
