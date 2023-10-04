package plugin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/bgrewell/gin-plugins/loader"
	"github.com/gin-gonic/gin"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/internal/config"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
)

func NewSubsystem(router *gin.Engine, log *zap.Logger, cfg *config.Configuration) *PluginSubsystem {
	ps := PluginSubsystem{
		Router:  router,
		Logger:  log.With(zap.String("module", "plugin")),
		Config:  cfg,
		enabled: cfg.Plugins.Enabled,
	}
	if ps.enabled {
		err := ps.Initialize()
		if err != nil {
			ps.Logger.Error("failed to initialize plugin subsystem", zap.Error(err))
			return nil
		}
	}
	return &ps
}

type PluginSubsystem struct {
	Router  *gin.Engine
	Logger  *zap.Logger
	Config  *config.Configuration
	enabled bool
}

func (ps *PluginSubsystem) Initialize() (err error) {

	group := &ps.Router.RouterGroup
	if ps.Config.Plugins.PluginGroup != "" {
		group = ps.Router.Group(ps.Config.Plugins.PluginGroup)
	}

	// Remap the plugin configs to use full path for key
	cm := make(map[string]*loader.PluginConfig)
	for k, v := range ps.Config.Plugins.Entries {

		// Deal with any poorly formed entries
		if v == nil {
			ps.Logger.Error("bad plugin entry", zap.String("name", k))
			continue
		}
		full := path.Join(ps.Config.Plugins.PluginDir, fmt.Sprintf("%s.plugin", k))
		v.PluginPath = full
		ps.Logger.Info("loaded configuration",
			zap.String("name", v.Name()),
			zap.Bool("enabled", v.Enabled),
			zap.String("path", v.PluginPath),
			zap.String("cookie", v.Cookie),
			zap.String("hash", v.Hash))

		if v.Hash != "" {
			ph, err := ComputeSHA256(v.PluginPath)
			if err != nil {
				ps.Logger.Error("failed to compute plugin hash",
					zap.Error(err),
					zap.String("name", v.Name()),
					zap.String("path", v.PluginPath))
			}
			if ph != v.Hash {
				ps.Logger.Warn("plugin not loaded hash check failed",
					zap.String("name", v.Name()),
					zap.String("path", v.PluginPath),
					zap.String("expected", v.Hash),
					zap.String("got", ph))
				continue
			}
		}
		cm[full] = v
	}

	l := loader.NewPluginLoader(ps.Config.Plugins.PluginDir, cm, group, ps.Config.Plugins.LoadUnconfigured)
	active, err := l.Initialize()
	if err != nil {
		return err
	}

	ps.Logger.Info("loaded plugins", zap.Int("count", len(active)))
	for idx, plug := range active {
		ps.Logger.Info("plugin activated",
			zap.Int("index", idx),
			zap.String("name", plug.Name),
			zap.String("path", plug.Path))
	}

	return nil
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
