package defaults

import (
	"time"

	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/pkg/structs"
	"github.com/Quad4-Software/cloud/services/gateway/pkg/config"
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
			Addr:   "127.0.0.1:9143",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9142",
			Namespace: "io.quad4.api",
			Protocol:  "tcp",
		},
		Service: config.Service{
			Name: "gateway",
		},
		Reva:                       shared.DefaultRevaConfig(),
		CommitShareToStorageGrant:  true,
		ShareFolder:                "Shares",
		DisableHomeCreationOnLogin: true,
		Cache: config.Cache{
			ProviderCacheStore:      "noop",
			ProviderCacheNodes:      []string{"127.0.0.1:9233"},
			ProviderCacheDatabase:   "cache-providers",
			ProviderCacheTTL:        300 * time.Second,
			CreateHomeCacheStore:    "memory",
			CreateHomeCacheNodes:    []string{"127.0.0.1:9233"},
			CreateHomeCacheDatabase: "cache-createhome",
			CreateHomeCacheTTL:      300 * time.Second,
		},

		FrontendPublicURL: "https://localhost:9200",

		AppRegistryEndpoint:       "io.quad4.api.app-registry",
		AuthAppEndpoint:           "io.quad4.api.auth-app",
		AuthBasicEndpoint:         "io.quad4.api.auth-basic",
		AuthGuestLinkEndpoint:     "io.quad4.api.auth-guest",
		AuthMachineEndpoint:       "io.quad4.api.auth-machine",
		AuthServiceEndpoint:       "io.quad4.api.auth-service",
		GroupsEndpoint:            "io.quad4.api.groups",
		PermissionsEndpoint:       "io.quad4.api.settings",
		SharingEndpoint:           "io.quad4.api.sharing",
		StoragePublicLinkEndpoint: "io.quad4.api.storage-publiclink",
		StorageSharesEndpoint:     "io.quad4.api.storage-shares",
		StorageUsersEndpoint:      "io.quad4.api.storage-users",
		UsersEndpoint:             "io.quad4.api.users",
		OCMEndpoint:               "io.quad4.api.ocm",

		StorageRegistry: config.StorageRegistry{
			Driver: "spaces",
			JSON:   "",
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

	if cfg.GRPC.TLS == nil && cfg.Commons != nil {
		cfg.GRPC.TLS = structs.CopyOrZeroValue(cfg.Commons.GRPCServiceTLS)
	}
}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
