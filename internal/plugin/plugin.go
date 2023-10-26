package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/bgrewell/gin-plugins/loader"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/interfaces"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/types/endpoint"
	"go.uber.org/zap"
)

// NewSubsystem creates a new instance of the Subsystem struct
func NewSubsystem(router *gin.Engine, log *zap.Logger, cfg *config.Configuration) interfaces.Subsystem {
	name := "plugin"
	ps := Subsystem{
		Router:  router,
		Logger:  log.With(zap.String("module", name)),
		Config:  cfg,
		enabled: cfg.Plugins.Enabled,
		name:    name,
	}
	ps.register()
	return &ps
}

// Subsystem handles plugin related functionalities
type Subsystem struct {
	Router    *gin.Engine
	Logger    *zap.Logger
	Config    *config.Configuration
	enabled   bool
	name      string // Subsystem name
	endpoints []endpoint.Endpoint
}

// register registers the routes that this module handles.
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	group := s.Config.Plugins.PluginGroup
	_ = group //TODO: clear warning

	// Remap the plugin configs to use full path for key
	cm := make(map[string]*loader.PluginConfig)
	for k, v := range s.Config.Plugins.Entries {

		// Deal with any poorly formed entries
		if v == nil {
			s.Logger.Error("bad plugin entry", zap.String("name", k))
			continue
		}
		full := path.Join(s.Config.Plugins.PluginDir, fmt.Sprintf("%s.plugin", k))
		v.PluginPath = full
		s.Logger.Info("loaded configuration",
			zap.String("name", v.Name()),
			zap.Bool("enabled", v.Enabled),
			zap.String("path", v.PluginPath),
			zap.String("cookie", v.Cookie),
			zap.String("hash", v.Hash))

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

	// TODO: Fix all of this, temporary hack to clear automated checks
	l := loader.NewPluginLoader(s.Config.Plugins.PluginDir, cm, &gin.Default().RouterGroup, s.Config.Plugins.LoadUnconfigured)
	active, err := l.Initialize()
	if err != nil {
		return
	}

	s.Logger.Info("loaded plugins", zap.Int("count", len(active)))
	for idx, plug := range active {
		s.Logger.Info("plugin activated",
			zap.Int("index", idx),
			zap.String("name", plug.Name),
			zap.String("path", plug.Path))
	}
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
func (s *Subsystem) Endpoints() []endpoint.Endpoint {
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
