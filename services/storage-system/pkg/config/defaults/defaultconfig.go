package defaults

import (
	"path/filepath"
	"time"

	"github.com/Quad4-Software/cloud/pkg/config/defaults"
	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/pkg/structs"
	"github.com/Quad4-Software/cloud/services/storage-system/pkg/config"
)

// FullDefaultConfig returns a fully initialized default configuration
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig returns a basic default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9217",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9215",
			Namespace: "io.quad4.api",
			Protocol:  "tcp",
		},
		HTTP: config.HTTPConfig{
			Addr:      "127.0.0.1:9216",
			Namespace: "io.quad4.web",
			Protocol:  "tcp",
		},
		Service: config.Service{
			Name: "storage-system",
		},
		Reva:          shared.DefaultRevaConfig(),
		DataServerURL: "http://localhost:9216/data",
		Driver:        "decomposed",
		Drivers: config.Drivers{
			Decomposed: config.DecomposedDriver{
				Root:                    filepath.Join(defaults.BaseDataPath(), "storage", "metadata"),
				MaxAcquireLockCycles:    20,
				LockCycleDurationFactor: 30,
			},
		},
		FileMetadataCache: config.Cache{
			Store:    "memory",
			Nodes:    []string{"127.0.0.1:9233"},
			Database: "storage-system",
			TTL:      24 * time.Hour,
		},
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "error"
	}
	if cfg.Reva == nil && cfg.Commons != nil {
		cfg.Reva = structs.CopyOrZeroValue(cfg.Commons.Reva)
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}

	if cfg.SystemUserAPIKey == "" && cfg.Commons != nil && cfg.Commons.SystemUserAPIKey != "" {
		cfg.SystemUserAPIKey = cfg.Commons.SystemUserAPIKey
	}

	if cfg.SystemUserID == "" && cfg.Commons != nil && cfg.Commons.SystemUserID != "" {
		cfg.SystemUserID = cfg.Commons.SystemUserID
	}

	if cfg.GRPC.TLS == nil && cfg.Commons != nil {
		cfg.GRPC.TLS = structs.CopyOrZeroValue(cfg.Commons.GRPCServiceTLS)
	}

}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
