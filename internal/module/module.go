package module

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/bgrewell/dtac-agent/internal/basic"
	"github.com/bgrewell/dtac-agent/internal/config"
	"github.com/bgrewell/dtac-agent/internal/interfaces"
	"github.com/bgrewell/dtac-agent/pkg/endpoint"
	"github.com/bgrewell/dtac-agent/pkg/modules"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

// NewSubsystem creates a new instance of the Subsystem struct
func NewSubsystem(log *zap.Logger, cfg *config.Configuration, tls *map[string]basic.TLSInfo) interfaces.Subsystem {
	name := "module"
	ms := Subsystem{
		Logger:  log.With(zap.String("module", name)),
		Config:  cfg,
		tls:     tls,
		enabled: cfg.Modules.Enabled,
		name:    name,
	}
	ms.register()
	return &ms
}

// Subsystem handles module related functionalities
type Subsystem struct {
	Logger  *zap.Logger
	Config  *config.Configuration
	tls     *map[string]basic.TLSInfo
	enabled bool
	name    string // Subsystem name
}

// register registers the routes that this module handles.
func (s *Subsystem) register() {
	if !s.Enabled() {
		s.Logger.Info("subsystem is disabled", zap.String("subsystem", s.Name()))
		return
	}

	group := s.Config.Modules.ModuleGroup

	// Remap the module configs to use full path for key
	cm := make(map[string]*modules.ModuleConfig)
	for k, v := range s.Config.Modules.Entries {

		// Deal with any poorly formed entries
		if v == nil {
			s.Logger.Error("bad module entry", zap.String("name", k))
			continue
		}
		full := path.Join(s.Config.Modules.ModuleDir, fmt.Sprintf("%s.module", k))
		if runtime.GOOS == "windows" {
			full = strings.Replace(full, "/", "", -1)
			full += ".exe"
		} else if runtime.GOOS == "darwin" {
			full += ".app"
		}

		v.ModulePath = full
		v.RootPath = group
		s.Logger.Info("loaded configuration",
			zap.String("name", v.Name()),
			zap.Bool("enabled", v.Enabled),
			zap.String("path", v.ModulePath),
			zap.String("hash", v.Hash),
			zap.String("root", v.RootPath),
			zap.String("config_key", full))

		if v.Hash != "" {
			mh, err := ComputeSHA256(v.ModulePath)
			if err != nil {
				s.Logger.Error("failed to compute module hash",
					zap.Error(err),
					zap.String("name", v.Name()),
					zap.String("path", v.ModulePath))
			}
			if mh != v.Hash {
				s.Logger.Warn("module not loaded hash check failed",
					zap.String("name", v.Name()),
					zap.String("path", v.ModulePath),
					zap.String("expected", v.Hash),
					zap.String("got", mh))
				continue
			}
		}
		cm[full] = v
	}

	// Check for TLS config
	var tlsKey, tlsCert, tlsCACert *string
	if s.Config.Modules.TLS.Enabled {
		profileName := s.Config.Modules.TLS.Profile
		if profile, ok := (*s.tls)[profileName]; ok {
			tlsCert = &profile.CertFilename
			tlsKey = &profile.KeyFilename
			tlsCACert = &profile.CAFilename
		}
	}

	loader := modules.NewModuleLoader(s.Config.Modules.ModuleDir, group, cm, s.Config.Modules.LoadUnconfigured, tlsCert, tlsKey, tlsCACert, s.Logger)
	active, err := loader.Initialize(s.Config.Auth.DefaultSecure)
	if err != nil {
		s.Logger.Error("failed to initialize modules", zap.Error(err))
		return
	}

	s.Logger.Info("loaded modules", zap.Int("count", len(active)))
	for idx, mod := range active {
		s.Logger.Info("module activated",
			zap.Int("index", idx),
			zap.String("name", mod.Name),
			zap.String("type", mod.ModuleType),
			zap.Strings("capabilities", mod.Capabilities),
			zap.String("path", mod.Path))
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

// Endpoints returns the endpoints for the subsystem (modules don't expose endpoints directly)
func (s *Subsystem) Endpoints() []*endpoint.Endpoint {
	// Modules don't expose REST/gRPC endpoints like plugins do
	// They may host web servers or other services independently
	return nil
}

// ComputeSHA256 computes the SHA256 hash of a file
func ComputeSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
