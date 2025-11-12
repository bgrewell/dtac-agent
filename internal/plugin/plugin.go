package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/bgrewell/dtac-agent/internal/basic"
	"github.com/bgrewell/dtac-agent/internal/config"
	"github.com/bgrewell/dtac-agent/internal/interfaces"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/plugins"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

// NewSubsystem creates a new instance of the Subsystem struct
func NewSubsystem(log *zap.Logger, cfg *config.Configuration, tls *map[string]basic.TLSInfo) interfaces.Subsystem {
	name := "plugin"
	ps := Subsystem{
		Logger:  log.With(zap.String("module", name)),
		Config:  cfg,
		tls:     tls,
		enabled: cfg.Plugins.Enabled,
		name:    name,
	}
	ps.register()
	return &ps
}

// Subsystem handles plugin related functionalities
type Subsystem struct {
	Logger    *zap.Logger
	Config    *config.Configuration
	tls       *map[string]basic.TLSInfo
	enabled   bool
	name      string // Subsystem name
	endpoints []*endpoint.Endpoint
}

// register registers the routes that this module handles.
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	group := s.Config.Plugins.PluginGroup

	// Remap the plugin configs to use full path for key
	cm := make(map[string]*plugins.PluginConfig)
	for k, v := range s.Config.Plugins.Entries {

		// Deal with any poorly formed entries
		if v == nil {
			s.Logger.Error("bad plugin entry", zap.String("name", k))
			continue
		}
		full := path.Join(s.Config.Plugins.PluginDir, fmt.Sprintf("%s.plugin", k))
		if runtime.GOOS == "windows" {
			full = strings.Replace(full, "/", "", -1)
			full += ".exe"
		} else if runtime.GOOS == "darwin" {
			full += ".app"
		}

		v.PluginPath = full
		v.RootPath = group
		s.Logger.Info("loaded configuration",
			zap.String("name", v.Name()),
			zap.Bool("enabled", v.Enabled),
			zap.String("path", v.PluginPath),
			zap.String("hash", v.Hash),
			zap.String("root", v.RootPath),
			zap.String("config_key", full))

		if v.Hash != "" {
			ph, err := ComputeSHA256(v.PluginPath)
			if err != nil {
				s.Logger.Error("failed to compute plugin hash",
					zap.Error(err),
					zap.String("name", v.Name()),
					zap.String("path", v.PluginPath))
			}
			if ph != v.Hash {
				s.Logger.Warn("plugin not loaded hash check failed",
					zap.String("name", v.Name()),
					zap.String("path", v.PluginPath),
					zap.String("expected", v.Hash),
					zap.String("got", ph))
				continue
			}
		}
		cm[full] = v
	}

	// Check for TLS config
	var tlsKey, tlsCert, tlsCACert *string
	if s.Config.Plugins.TLS.Enabled {
		profileName := s.Config.Plugins.TLS.Profile
		if profile, ok := (*s.tls)[profileName]; ok {
			tlsCert = &profile.CertFilename
			tlsKey = &profile.KeyFilename
			tlsCACert = &profile.CAFilename
		}
	}

	loader := plugins.NewPluginLoader(s.Config.Plugins.PluginDir, group, cm, s.Config.Plugins.LoadUnconfigured, tlsCert, tlsKey, tlsCACert, s.Logger)
	active, err := loader.Initialize(s.Config.Auth.DefaultSecure)
	if err != nil {
		s.Logger.Error("failed to initialize plugins", zap.Error(err))
		return
	}

	s.Logger.Info("loaded plugins", zap.Int("count", len(active)))
	for idx, plug := range active {
		s.Logger.Info("plugin activated",
			zap.Int("index", idx),
			zap.String("name", plug.Name),
			zap.String("path", plug.Path))
	}

	// Manipulate the endpoints
	for _, ep := range loader.Endpoints() {
		// Intentionally shadow 'ep' so the closure captures the correct value
		ep := ep

		// If plugins have a group namespace then append it
		if group != "" {
			ep.Path = path.Join(group, ep.Path)
		}

		// Ensure all functions point to the plugin loaders shim
		ep.Function = func(in *endpoint.Request) (out *endpoint.Response, err error) {
			return loader.CallShim(ep, in)
		}
	}

	s.endpoints = loader.Endpoints()
}

// Enabled returns true if the subsystem is enabled
func (s *Subsystem) Enabled() bool {
	return s.enabled
}

// Name returns the name of the subsystem
func (s *Subsystem) Name() string {
	return s.name
}

// Endpoints returns an array of endpoints that this Subsystem handles
func (s *Subsystem) Endpoints() []*endpoint.Endpoint {
	return s.endpoints
}

// ComputeSHA256 computes the SHA-256 hash of a file specified by its path.
func ComputeSHA256(filePath string) (string, error) {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new hash calculator
	hasher := sha256.New()

	// Read the file content and compute the hash
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", err
	}

	// Convert the hash result to a hex string
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash, nil
}
