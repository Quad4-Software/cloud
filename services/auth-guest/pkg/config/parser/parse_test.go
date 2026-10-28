// Copyright 2026 Quad4 Cloud GmbH <mail@quad4.io>
// SPDX-License-Identifier: Apache-2.0

package parser

import (
	"testing"

	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/config"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		jwtSecret string
		wantErr   bool
	}{
		{name: "jwt secret set", jwtSecret: "reva-secret"},
		{name: "missing jwt secret", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				TokenManager: &config.TokenManager{JWTSecret: tt.jwtSecret},
			}

			err := Validate(cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
