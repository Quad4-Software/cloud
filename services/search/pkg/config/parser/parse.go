package parser

import (
	"errors"
	"fmt"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"

	occfg "github.com/Quad4-Software/cloud/pkg/config"
	"github.com/Quad4-Software/cloud/pkg/kql"
	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/services/search/pkg/config"
	"github.com/Quad4-Software/cloud/services/search/pkg/config/defaults"
	"github.com/Quad4-Software/cloud/services/search/pkg/search"

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

	return Validate(cfg)
}

func Validate(cfg *config.Config) error {
	if cfg.Events.Disabled && cfg.GRPC.Disabled {
		return shared.AllComponentsDisabledError(cfg.Service.Name)
	}

	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}

	if cfg.ServiceAccount.ServiceAccountID == "" {
		return shared.MissingServiceAccountID(cfg.Service.Name)
	}
	if cfg.ServiceAccount.ServiceAccountSecret == "" {
		return shared.MissingServiceAccountSecret(cfg.Service.Name)
	}

	return validateQueryTemplate(cfg.QueryTemplate)
}

func validateQueryTemplate(template string) error {
	if template == "" {
		return nil
	}

	if strings.Count(template, config.QueryPlaceholder) != 1 {
		return errors.New("search query template: must contain " + config.QueryPlaceholder + " exactly once")
	}
	wrapped := "(" + config.QueryPlaceholder + ")"
	if !strings.Contains(template, wrapped) {
		return errors.New("search query template: " + config.QueryPlaceholder + " must be wrapped as " + wrapped)
	}

	expanded := search.ExpandQueryTemplate(template, "a", &provider.StorageSpace{})

	if strings.Contains(expanded, "{{") || strings.Contains(expanded, "}}") {
		return errors.New("search query template: unknown placeholder")
	}

	if _, err := (kql.Builder{}).Build(expanded); err != nil {
		return fmt.Errorf("search query template: %w", err)
	}

	return nil
}
