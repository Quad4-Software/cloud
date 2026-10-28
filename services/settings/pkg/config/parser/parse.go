package parser

import (
	"errors"

	occfg "github.com/Quad4-Software/cloud/pkg/config"
	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/services/settings/pkg/config"
	"github.com/Quad4-Software/cloud/services/settings/pkg/config/defaults"

	"github.com/Quad4-Software/cloud/pkg/config/envdecode"
)

// ParseConfig loads configuration from known paths.
func ParseConfig(cfg *config.Config) error {
	err := occfg.BindSourcesToStructs(cfg.Service.Name, cfg)
	if err != nil {
		return err
	}

	defaults.EnsureDefaults(cfg)
	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	defaults.Sanitize(cfg)

	if err := defaults.LoadBundles(cfg); err != nil {
		return err
	}

	return Validate(cfg)
}

func Validate(cfg *config.Config) error {
	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}

	if cfg.Metadata.SystemUserAPIKey == "" {
		return shared.MissingSystemUserApiKeyError(cfg.Service.Name)
	}

	if cfg.SetupDefaultAssignments && cfg.AdminUserID == "" {
		return shared.MissingAdminUserID(cfg.Service.Name)
	}

	if len(cfg.ServiceAccountIDs) == 0 {
		return shared.MissingServiceAccountID(cfg.Service.Name)
	}

	return nil
}
