// Copyright 2026 Quad4 Cloud GmbH <mail@quad4.io>
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/config"
	"github.com/spf13/cobra"
)

// Health is the entrypoint for the health command.
func Health(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "Check health status",
		RunE: func(cmd *cobra.Command, args []string) error {
			// not implemented
			return nil
		},
	}
}
