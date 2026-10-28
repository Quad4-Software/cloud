package debug

import (
	"net/http"
	"strconv"

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
		options.Config.Nats.EnableTLS,
		options.Config.Nats.TLSSkipVerifyClientCert,
		options.Config.Nats.TLSCert,
	)
	// For nats readiness and liveness checks are identical
	// the nats server will neither be healthy nor ready when it can not reach the nats server/cluster
	checkHandler := handlers.NewCheckHandler(
		handlers.NewCheckHandlerConfiguration().
			WithLogger(options.Logger).
			WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Nats.Host+":"+strconv.Itoa(options.Config.Nats.Port), secureOption)),
	)

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(checkHandler),
		debug.Ready(checkHandler),
	), nil
}
