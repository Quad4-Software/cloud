package debug

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/dutchcoders/go-clamd"

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
		options.Config.Events.TLSRootCACertificate,
	)
	readyHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("nats reachability", checks.NewNatsCheck(options.Config.Events.Endpoint, secureOption, nats.Auth(options.Config.Events.AuthUsername, options.Config.Events.AuthPassword))).
		WithCheck("antivirus reachability", func(ctx context.Context) error {
			cfg := options.Config
			switch cfg.Scanner.Type {
			default:
				return errors.New("no antivirus configured")
			case "clamav":
				return clamd.NewClamd(cfg.Scanner.ClamAV.Socket).Ping()
			case "icap":
				u, err := url.Parse(cfg.Scanner.ICAP.URL)
				if err != nil {
					return err
				}
				return checks.NewTCPCheck(u.Host)(ctx)
			}
		})

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
