// Copyright 2026 Quad4 Cloud GmbH <mail@quad4.io>
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"os"

	"github.com/Quad4-Software/cloud/pkg/clihelper"
	"github.com/Quad4-Software/cloud/services/auth-guest/pkg/config"
	"github.com/spf13/cobra"
)

// GetCommands provides all commands for this service
func GetCommands(cfg *config.Config) []*cobra.Command {
	return []*cobra.Command{
		// start this service
		Server(cfg),

		// interaction with this service

		// infos about this service
		Health(cfg),
		Version(cfg),
	}
}

// Execute is the entry point for the auth-guest command.
func Execute(cfg *config.Config) error {
	app := clihelper.DefaultApp(&cobra.Command{
		Use:   "auth-guest",
		Short: "starts auth-guest service",
	})
	app.AddCommand(GetCommands(cfg)...)
	app.SetArgs(os.Args[1:])
	return app.ExecuteContext(cfg.Context)
}
