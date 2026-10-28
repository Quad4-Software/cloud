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
		options.Config.Postprocessing.Events.EnableTLS,
		options.Config.Postprocessing.Events.TLSInsecure,
		options.Config.Postprocessing.Events.TLSRootCACertificate,
	)
	readyHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Postprocessing.Events.Endpoint, secureOption, nats.Auth(options.Config.Postprocessing.Events.AuthUsername, options.Config.Postprocessing.Events.AuthPassword)))

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Ready(handlers.NewCheckHandler(readyHandlerConfiguration)),
	), nil
}
