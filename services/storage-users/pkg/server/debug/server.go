package debug

import (
	"net/http"

	"github.com/Quad4-Software/cloud/pkg/checks"
	"github.com/Quad4-Software/cloud/pkg/handlers"
	"github.com/Quad4-Software/cloud/pkg/nats"
	"github.com/Quad4-Software/cloud/pkg/service/debug"
	"github.com/Quad4-Software/cloud/pkg/version"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	secureOption := nats.Secure(
		options.Config.Events.EnableTLS,
		options.Config.Events.TLSInsecure,
		options.Config.Events.TLSRootCaCertPath,
	)
	readyHandler := handlers.NewCheckHandler(handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Events.Addr, secureOption, nats.Auth(options.Config.Events.AuthUsername, options.Config.Events.AuthPassword))).
		WithCheck("grpc reachability", checks.NewGRPCCheck(options.Config.GRPC.Addr)),
	)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Context(options.Context),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Ready(readyHandler),
		//debug.CorsAllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
		//debug.CorsAllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
		//debug.CorsAllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
		//debug.CorsAllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
	), nil
}
