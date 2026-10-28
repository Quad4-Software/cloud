package command

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/Quad4-Software/cloud/pkg/config/configlog"
	ocEvents "github.com/Quad4-Software/cloud/pkg/events"
	"github.com/Quad4-Software/cloud/pkg/generators"
	"github.com/Quad4-Software/cloud/pkg/log"
	"github.com/Quad4-Software/cloud/pkg/registry"
	"github.com/Quad4-Software/cloud/pkg/runner"
	ogrpc "github.com/Quad4-Software/cloud/pkg/service/grpc"
	"github.com/Quad4-Software/cloud/pkg/tracing"
	"github.com/Quad4-Software/cloud/pkg/version"
	ehsvc "github.com/Quad4-Software/cloud/protogen/gen/quad4/services/eventhistory/v0"
	settingssvc "github.com/Quad4-Software/cloud/protogen/gen/quad4/services/settings/v0"
	"github.com/Quad4-Software/cloud/services/userlog/pkg/config"
	"github.com/Quad4-Software/cloud/services/userlog/pkg/config/parser"
	"github.com/Quad4-Software/cloud/services/userlog/pkg/metrics"
	"github.com/Quad4-Software/cloud/services/userlog/pkg/server/debug"
	"github.com/Quad4-Software/cloud/services/userlog/pkg/server/http"

	"github.com/opencloud-eu/reva/v2/pkg/events"
	"github.com/opencloud-eu/reva/v2/pkg/events/stream"
	"github.com/opencloud-eu/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/opencloud-eu/reva/v2/pkg/store"
	"github.com/spf13/cobra"
	microstore "go-micro.dev/v4/store"
)

// all events we care about
var _registeredEvents = []events.Unmarshaller{
	// file related
	events.PostprocessingStepFinished{},

	// space related
	events.SpaceDisabled{},
	events.SpaceDeleted{},
	events.SpaceShared{},
	events.SpaceUnshared{},
	events.SpaceMembershipExpired{},

	// share related
	events.ShareCreated{},
	events.ShareRemoved{},
	events.ShareExpired{},

	// misc
	ocEvents.ResourceMention{},
}

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
			tracerProvider, err := tracing.GetTraceProvider(cmd.Context(), cfg.Commons.TracesExporter, cfg.Service.Name)
			if err != nil {
				return err
			}

			grpcClient, err := ogrpc.NewClient(
				append(ogrpc.GetClientOptions(cfg.GRPCClientTLS), ogrpc.WithTraceProvider(tracerProvider))...,
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

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			connName := generators.GenerateConnectionName(cfg.Service.Name, generators.NTypeBus)
			stream, err := stream.NatsFromConfig(connName, false, stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			st := store.Create(
				store.Store(cfg.Persistence.Store),
				store.TTL(cfg.Persistence.TTL),
				microstore.Nodes(cfg.Persistence.Nodes...),
				microstore.Database(cfg.Persistence.Database),
				microstore.Table(cfg.Persistence.Table),
				store.Authentication(cfg.Persistence.AuthUsername, cfg.Persistence.AuthPassword),
				store.TLSEnabled(cfg.Persistence.EnableTLS),
				store.TLSInsecure(cfg.Persistence.TLSInsecure),
				store.TLSRootCA(cfg.Persistence.TLSRootCACertificate),
			)

			tm, err := pool.StringToTLSMode(cfg.GRPCClientTLS.Mode)
			if err != nil {
				return err
			}
			gatewaySelector, err := pool.GatewaySelector(
				cfg.RevaGateway,
				pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
				pool.WithTLSMode(tm),
				pool.WithRegistry(registry.GetRegistry()),
				pool.WithTracerProvider(tracerProvider),
			)
			if err != nil {
				return fmt.Errorf("could not get reva client selector: %s", err)
			}

			hClient := ehsvc.NewEventHistoryService("io.quad4.api.eventhistory", grpcClient)
			vClient := settingssvc.NewValueService("io.quad4.api.settings", grpcClient)
			rClient := settingssvc.NewRoleService("io.quad4.api.settings", grpcClient)

			gr := runner.NewGroup()
			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
					http.Store(st),
					http.Stream(stream),
					http.GatewaySelector(gatewaySelector),
					http.History(hClient),
					http.Value(vClient),
					http.Role(rClient),
					http.RegisteredEvents(_registeredEvents),
					http.TracerProvider(tracerProvider),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
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
