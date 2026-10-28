package debug

import (
	"encoding/json"
	"net/http"

	"github.com/ggwhite/go-masker"

	"github.com/Quad4-Software/cloud/pkg/checks"
	"github.com/Quad4-Software/cloud/pkg/handlers"
	"github.com/Quad4-Software/cloud/pkg/nats"
	"github.com/Quad4-Software/cloud/pkg/service/debug"
	"github.com/Quad4-Software/cloud/pkg/version"
	"github.com/Quad4-Software/cloud/services/proxy/pkg/config"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	healthHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("web reachability", checks.NewHTTPCheck(options.Config.HTTP.Addr))

	secureOption := nats.Secure(
		options.Config.Events.EnableTLS,
		options.Config.Events.TLSInsecure,
		options.Config.Events.TLSRootCACertificate,
	)
	readyHandlerConfiguration := healthHandlerConfiguration.
		WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Events.Endpoint, secureOption, nats.Auth(options.Config.Events.AuthUsername, options.Config.Events.AuthPassword)))

	var configDumpFunc http.HandlerFunc = configDump(options.Config)
	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Health(handlers.NewCheckHandler(healthHandlerConfiguration)),
		debug.Ready(handlers.NewCheckHandler(readyHandlerConfiguration)),
		debug.ConfigDump(configDumpFunc),
	), nil
}

// configDump implements the config dump
func configDump(cfg *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		maskedCfg, err := masker.Struct(cfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		b, err := json.Marshal(maskedCfg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		_, _ = w.Write(b)
	}
}
