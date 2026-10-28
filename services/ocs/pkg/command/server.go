package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/Quad4-Software/cloud/pkg/config/configlog"
	"github.com/Quad4-Software/cloud/pkg/log"
	"github.com/Quad4-Software/cloud/pkg/runner"
	ogrpc "github.com/Quad4-Software/cloud/pkg/service/grpc"
	"github.com/Quad4-Software/cloud/pkg/tracing"
	"github.com/Quad4-Software/cloud/pkg/version"
	"github.com/Quad4-Software/cloud/services/ocs/pkg/config"
	"github.com/Quad4-Software/cloud/services/ocs/pkg/config/parser"
	"github.com/Quad4-Software/cloud/services/ocs/pkg/metrics"
	"github.com/Quad4-Software/cloud/services/ocs/pkg/server/debug"
	"github.com/Quad4-Software/cloud/services/ocs/pkg/server/http"

	"github.com/spf13/cobra"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.Configure(cfg.Service.Name, cfg.Commons, cfg.LogLevel)
			traceProvider, err := tracing.GetTraceProvider(cmd.Context(), cfg.Commons.TracesExporter, cfg.Service.Name)
			if err != nil {
				return err
			}

			cfg.GrpcClient, err = ogrpc.NewClient(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			if err != nil {
				return err
			}

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			metrics := metrics.New()
			metrics.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			gr := runner.NewGroup()
			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(metrics),
					http.TraceProvider(traceProvider),
				)

				if err != nil {
					logger.Error().
						Err(err).
						Str("transport", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(runner.NewGoMicroHttpServerRunner(cfg.Service.Name+".http", server))
			}

			{
				debugServer, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(runner.NewGolangHttpServerRunner(cfg.Service.Name+".debug", debugServer))
			}

			grResults := gr.Run(ctx)

			// return the first non-nil error found in the results
			for _, grResult := range grResults {
				if grResult.RunnerError != nil {
					return grResult.RunnerError
				}
			}
			return nil
		},
	}
}
