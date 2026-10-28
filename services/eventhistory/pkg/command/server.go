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
	ehsvc "github.com/Quad4-Software/cloud/protogen/gen/quad4/services/eventhistory/v0"
	"github.com/Quad4-Software/cloud/services/eventhistory/pkg/config"
	"github.com/Quad4-Software/cloud/services/eventhistory/pkg/config/parser"
	"github.com/Quad4-Software/cloud/services/eventhistory/pkg/metrics"
	"github.com/Quad4-Software/cloud/services/eventhistory/pkg/server/consumer"
	"github.com/Quad4-Software/cloud/services/eventhistory/pkg/server/debug"
	"github.com/Quad4-Software/cloud/services/eventhistory/pkg/server/grpc"
	svc "github.com/Quad4-Software/cloud/services/eventhistory/pkg/service"
	"github.com/opencloud-eu/reva/v2/pkg/store"

	"github.com/spf13/cobra"
	microstore "go-micro.dev/v4/store"
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

			cfg.GrpcClient, err = ogrpc.NewClient(
				append(ogrpc.GetClientOptions(cfg.GRPCClientTLS), ogrpc.WithTraceProvider(traceProvider))...,
			)
			if err != nil {
				return err
			}

			var cancel context.CancelFunc
			if cfg.Context == nil {
				cfg.Context, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}
			ctx := cfg.Context

			m := metrics.New()
			m.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			gr := runner.NewGroup()

			st := store.Create(
				store.Store(cfg.Store.Store),
				store.TTL(cfg.Store.TTL),
				microstore.Nodes(cfg.Store.Nodes...),
				microstore.Database(cfg.Store.Database),
				microstore.Table(cfg.Store.Table),
				store.Authentication(cfg.Store.AuthUsername, cfg.Store.AuthPassword),
				store.TLSEnabled(cfg.Store.EnableTLS),
				store.TLSInsecure(cfg.Store.TLSInsecure),
				store.TLSRootCA(cfg.Store.TLSRootCACertificate),
			)

			if !cfg.Events.Disabled {
				_, err = consumer.NewConsumer(
					consumer.Logger(logger),
					consumer.Config(cfg),
					consumer.Persistence(st),
				)
				if err != nil {
					return err
				}
			} else {
				logger.Info().Msg("event listening disabled, not starting event consumer")
			}

			if !cfg.GRPC.Disabled {
				eh, err := svc.NewEventHistoryService(cfg, st, logger)
				if err != nil {
					return err
				}

				service := grpc.NewService(
					grpc.Logger(logger),
					grpc.Context(ctx),
					grpc.Config(cfg),
					grpc.Name(cfg.Service.Name),
					grpc.Namespace(cfg.GRPC.Namespace),
					grpc.Address(cfg.GRPC.Addr),
					grpc.Metrics(m),
					grpc.TraceProvider(traceProvider),
				)

				if err := ehsvc.RegisterEventHistoryServiceHandler(service.Server(), eh); err != nil {
					return err
				}

				gr.Add(runner.NewGoMicroGrpcServerRunner(cfg.Service.Name+".grpc", service))
			} else {
				logger.Info().Msg("gRPC server disabled, not starting gRPC service")
			}

			debugServer, err := debug.Server(
				debug.Logger(logger),
				debug.Context(ctx),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Info().Err(err).Str("server", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(runner.NewGolangHttpServerRunner(cfg.Service.Name+".debug", debugServer))

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
