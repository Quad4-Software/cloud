// Copyright 2026 Quad4 Cloud GmbH <mail@quad4.io>
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"errors"

	occfg "github.com/Quad4-Software/cloud/pkg/config"
	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/config"
	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/config/defaults"

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

// Validate validates the config
func Validate(cfg *config.Config) error {
	if cfg.TokenManager == nil || cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}
	return nil
}
