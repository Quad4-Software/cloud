// Copyright 2026 Quad4 Cloud GmbH <mail@quad4.io>
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/config"
	"github.com/spf13/cobra"
)

// Version prints the service versions of all running instances.
func Version(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of this binary and the running service instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			// not implemented
			return nil
		},
	}
}
