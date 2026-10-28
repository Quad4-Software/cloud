package command

import (
	"github.com/Quad4-Software/cloud/quad4/pkg/register"
	"github.com/Quad4-Software/cloud/quad4/pkg/runtime"
	"github.com/Quad4-Software/cloud/pkg/config"
	"github.com/Quad4-Software/cloud/pkg/config/configlog"
	"github.com/Quad4-Software/cloud/pkg/config/parser"

	"github.com/spf13/cobra"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "start a fullstack server (runtime and all services in supervised mode)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return configlog.ReturnError(parser.ParseConfig(cfg, false))
		},
		GroupID: CommandGroupServer,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Prefer the in-memory registry as the default when running in single-binary mode
			r := runtime.New(cfg)
			return r.Start(cmd.Context())
		},
	}
}

func init() {
	register.AddCommand(Server)
}
