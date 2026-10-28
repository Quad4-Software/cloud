package defaults

import (
	"strings"

	"github.com/Quad4-Software/cloud/services/webfinger/pkg/config"
	"github.com/Quad4-Software/cloud/services/webfinger/pkg/relations"
)

var (
	nativeAppScopes = []string{"openid", "profile", "email", "offline_access"}
	webAppScopes    = []string{"openid", "profile", "email"}
)

// FullDefaultConfig returns a fully initialized default configuration
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig returns a basic default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9279",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9275",
			Root:      "/",
			Namespace: "io.quad4.web",
			CORS: config.CORS{
				AllowedOrigins:   []string{"https://localhost:9200"},
				AllowCredentials: false,
			},
		},
		Service: config.Service{
			Name: "webfinger",
		},

		Quad4CloudURL: "https://localhost:9200",
		Relations:    []string{relations.OpenIDConnectRel, relations.Quad4CloudInstanceRel},
		Instances: []config.Instance{
			{
				Claim: "sub",
				Regex: ".+",
				Href:  "{{.OC_URL}}",
				Titles: map[string]string{
					"en": "Quad4 Cloud Instance",
				},
			},
		},
		IDP:                 "https://localhost:9200",
		Insecure:            false,
		AndroidClientID:     "Quad4CloudAndroid",
		AndroidClientScopes: nativeAppScopes,
		DesktopClientID:     "Quad4CloudDesktop",
		DesktopClientScopes: nativeAppScopes,
		IOSClientID:         "Quad4CloudIOS",
		IOSClientScopes:     nativeAppScopes,
		WebClientID:         "web",
		WebClientScopes:     webAppScopes,
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "error"
	}

	if cfg.Commons != nil {
		cfg.HTTP.TLS = cfg.Commons.HTTPServiceTLS
	}

	if (cfg.Commons != nil && cfg.Commons.Quad4CloudURL != "") &&
		(cfg.HTTP.CORS.AllowedOrigins == nil ||
			len(cfg.HTTP.CORS.AllowedOrigins) == 1 &&
				cfg.HTTP.CORS.AllowedOrigins[0] == "https://localhost:9200") {
		cfg.HTTP.CORS.AllowedOrigins = []string{cfg.Commons.Quad4CloudURL}
	}
}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}

	cfg.OIDCClientConfigs = map[string]config.OIDCClientConfig{
		"android": {
			ClientID: cfg.AndroidClientID,
			Scopes:   cfg.AndroidClientScopes,
		},
		"desktop": {
			ClientID: cfg.DesktopClientID,
			Scopes:   cfg.DesktopClientScopes,
		},
		"ios": {
			ClientID: cfg.IOSClientID,
			Scopes:   cfg.IOSClientScopes,
		},
		"web": {
			ClientID: cfg.WebClientID,
			Scopes:   cfg.WebClientScopes,
		},
	}
}
