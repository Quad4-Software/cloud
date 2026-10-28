package defaults

import (
	"strings"
	"time"

	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/pkg/structs"
	"github.com/Quad4-Software/cloud/services/userlog/pkg/config"
)

// FullDefaultConfig returns the full default config
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig return the default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9214",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		Service: config.Service{
			Name: "userlog",
		},
		Events: config.Events{
			Endpoint:  "127.0.0.1:9233",
			Cluster:   "quad4-cluster",
			EnableTLS: false,
		},
		MaxConcurrency: 1,
		Persistence: config.Persistence{
			Store:    "memory",
			Database: "userlog",
			Table:    "events",
			TTL:      time.Hour * 336,
		},
		RevaGateway: shared.DefaultRevaConfig().Address,
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9210",
			Root:      "/",
			Namespace: "io.quad4.web",
			CORS: config.CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With", "X-Request-Id", "Ocs-Apirequest"},
				AllowCredentials: true,
			},
		},
	}
}

// EnsureDefaults ensures the config contains default values
func EnsureDefaults(cfg *config.Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "error"
	}

	if cfg.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}

	if cfg.Commons != nil {
		cfg.HTTP.TLS = cfg.Commons.HTTPServiceTLS
	}

}

// Sanitize sanitizes the config
func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
	if cfg.MaxConcurrency < 1 {
		cfg.MaxConcurrency = 5
	}
}
