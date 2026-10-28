package defaults

import (
	"time"

	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/pkg/structs"
	"github.com/Quad4-Software/cloud/services/notifications/pkg/config"
)

// FullDefaultConfig returns a fully initialized default configuration
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// NOTE: Most of this configuration is not needed to keep it as simple as possible
// TODO: Clean up unneeded configuration

// DefaultConfig returns a basic default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9174",
			Zpages: false,
			Pprof:  false,
		},
		Service: config.Service{
			Name: "notifications",
		},
		WebUIURL: "https://localhost:9200",
		Notifications: config.Notifications{
			SMTP: config.SMTP{
				Encryption: "none",
			},
			Events: config.Events{
				Endpoint:  "127.0.0.1:9233",
				Cluster:   "quad4-cluster",
				EnableTLS: false,
			},
			RevaGateway: shared.DefaultRevaConfig().Address,
		},
		Store: config.Store{
			Store:    "nats-js-kv",
			Nodes:    []string{"127.0.0.1:9233"},
			Database: "notifications",
			Table:    "",
			TTL:      336 * time.Hour,
		},
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "error"
	}

	if cfg.Notifications.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.Notifications.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}
}

// Sanitize sanitizes the configuration
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
