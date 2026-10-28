package defaults

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Quad4-Software/cloud/pkg/config/defaults"
	"github.com/Quad4-Software/cloud/pkg/shared"
	"github.com/Quad4-Software/cloud/pkg/structs"
	"github.com/Quad4-Software/cloud/services/proxy/pkg/config"
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
			Addr:  "127.0.0.1:9205",
			Token: "",
		},
		HTTP: config.HTTP{
			Addr:      "0.0.0.0:9200",
			Root:      "/",
			Namespace: "io.quad4.web",
			TLSCert:   path.Join(defaults.BaseDataPath(), "proxy", "server.crt"),
			TLSKey:    path.Join(defaults.BaseDataPath(), "proxy", "server.key"),
			TLS:       true,
			Client: config.Client{
				ForceAttemptHTTP2:     false,
				DialTimeout:           30 * time.Second,
				DialKeepAlive:         30 * time.Second,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
		Service: config.Service{
			Name: "proxy",
		},
		OIDC: config.OIDC{
			Issuer: "https://localhost:9200",

			AccessTokenVerifyMethod: config.AccessTokenVerificationJWT,
			SkipUserInfo:            false,
			UserinfoCache: &config.Cache{
				Store:    "memory",
				Nodes:    []string{"127.0.0.1:9233"},
				Database: "cache-userinfo",
				TTL:      time.Second * 10,
			},
			JWKS: config.JWKS{
				RefreshInterval:   60, // minutes
				RefreshRateLimit:  60, // seconds
				RefreshTimeout:    10, // seconds
				RefreshUnknownKID: true,
			},
		},
		PolicySelector: nil,
		RoleAssignment: config.RoleAssignment{
			Driver: "default",
			// this default is only relevant when Driver is set to "oidc"
			OIDCRoleMapper: config.OIDCRoleMapper{
				RoleClaim: "roles",
				RolesMap: []config.RoleMapping{
					{RoleName: "admin", ClaimValue: "quad4Admin"},
					{RoleName: "spaceadmin", ClaimValue: "quad4SpaceAdmin"},
					{RoleName: "user", ClaimValue: "quad4User"},
					{RoleName: "user-light", ClaimValue: "quad4Guest"},
				},
			},
		},
		Reva: shared.DefaultRevaConfig(),
		PreSignedURL: config.PreSignedURL{
			AllowedHTTPMethods: []string{"GET"},
			Enabled:            true,
			SigningKeys: &config.SigningKeys{
				Store:              "nats-js-kv", // signing keys are written by ocs, so we cannot use memory. It is not shared.
				Nodes:              []string{"127.0.0.1:9233"},
				TTL:                time.Hour * 12,
				DisablePersistence: true,
			},
		},
		TransferTimeout:       24 * time.Hour,
		AccountBackend:        "cs3",
		UserOIDCClaim:         "preferred_username",
		UserCS3Claim:          "username",
		AutoprovisionAccounts: false,
		AutoProvisionClaims: config.AutoProvisionClaims{
			Username:    "preferred_username",
			Email:       "email",
			DisplayName: "name",
			Groups:      "groups",
		},
		EnableBasicAuth:               false,
		InsecureBackends:              false,
		CSPConfigFileLocation:         "",
		CSPConfigFileOverrideLocation: "",
		Events: config.Events{
			Endpoint:  "127.0.0.1:9233",
			Cluster:   "quad4-cluster",
			EnableTLS: false,
		},
		AuthMiddleware: config.AuthMiddleware{
			AllowAppAuth: true,
		},
		GuestLinkAuth: config.GuestLinkAuth{
			CookieName: "__Host-oc_guest_session",
		},
	}
}

// DefaultPolicies returns the default proxy policies.
func DefaultPolicies() []config.Policy {
	return []config.Policy{
		{
			Name: "default",
			Routes: []config.Route{
				{
					Endpoint:    "/",
					Service:     "io.quad4.web.web",
					Unprotected: true,
				},
				{
					Endpoint:    "/.well-known/ocm",
					Service:     "io.quad4.web.ocm",
					Unprotected: true,
				},
				{
					Endpoint:    "/.well-known/webfinger",
					Service:     "io.quad4.web.webfinger",
					Unprotected: true,
				},
				{
					Endpoint:    "/.well-known/openid-configuration",
					Service:     "io.quad4.web.idp",
					Unprotected: true,
				},
				{
					Endpoint: "/branding/logo",
					Service:  "io.quad4.web.web",
				},
				{
					Endpoint: "/announcement",
					Service:  "io.quad4.web.web",
				},
				{
					Endpoint:    "/konnect/",
					Service:     "io.quad4.web.idp",
					Unprotected: true,
				},
				{
					Endpoint:    "/signin/",
					Service:     "io.quad4.web.idp",
					Unprotected: true,
				},
				{
					Endpoint: "/archiver",
					Service:  "io.quad4.web.frontend",
				},
				{
					// reroute oc10 notifications endpoint to userlog service
					Endpoint: "/ocs/v2.php/apps/notifications/api/v1/notifications/sse",
					Service:  "io.quad4.sse.sse",
				},
				{
					// reroute oc10 notifications endpoint to userlog service
					Endpoint: "/ocs/v2.php/apps/notifications/api/v1/notifications",
					Service:  "io.quad4.web.userlog",
				},
				{
					Type:     config.RegexRoute,
					Endpoint: "/ocs/v[12].php/cloud/user/signing-key", // only `user/signing-key` is left in quad4-ocs
					Service:  "io.quad4.web.ocs",
				},
				{
					Type:        config.RegexRoute,
					Endpoint:    "/ocs/v[12].php/config",
					Service:     "io.quad4.web.frontend",
					Unprotected: true,
				},
				// OCM WAYF public endpoints
				{
					Endpoint:    "/sciencemesh/federations",
					Service:     "io.quad4.web.ocm",
					Unprotected: true,
				},
				{
					Endpoint:    "/sciencemesh/discover",
					Service:     "io.quad4.web.ocm",
					Unprotected: true,
				},
				// General sciencemesh endpoints
				{
					Endpoint: "/sciencemesh/",
					Service:  "io.quad4.web.ocm",
				},
				{
					Endpoint: "/ocm/",
					Service:  "io.quad4.web.ocm",
				},
				{
					Endpoint: "/ocs/",
					Service:  "io.quad4.web.frontend",
				},
				{
					Type:     config.QueryRoute,
					Endpoint: "/remote.php/?preview=1",
					Service:  "io.quad4.web.webdav",
				},
				// TODO the actual REPORT goes to /dav/files/{username}, which is user specific ... how would this work in a spaces world?
				// TODO what paths are returned? the href contains the full path so it should be possible to return urls from other spaces?
				// TODO or we allow a REPORT on /dav/spaces to search all spaces and /dav/space/{spaceid} to search a specific space
				// send webdav REPORT requests to search service
				{
					Type:     config.RegexRoute,
					Method:   "REPORT",
					Endpoint: "(/remote.php)?/(web)?dav",
					Service:  "io.quad4.web.webdav",
				},
				{
					Type:     config.QueryRoute,
					Endpoint: "/dav/?preview=1",
					Service:  "io.quad4.web.webdav",
				},
				{
					Type:     config.QueryRoute,
					Endpoint: "/webdav/?preview=1",
					Service:  "io.quad4.web.webdav",
				},
				{
					Endpoint: "/remote.php/",
					Service:  "io.quad4.web.frontend",
				},
				{
					Endpoint: "/dav/",
					Service:  "io.quad4.web.frontend",
				},
				{
					Endpoint: "/webdav/",
					Service:  "io.quad4.web.frontend",
				},
				{
					Endpoint:    "/status",
					Service:     "io.quad4.web.frontend",
					Unprotected: true,
				},
				{
					Endpoint:    "/status.php",
					Service:     "io.quad4.web.frontend",
					Unprotected: true,
				},
				{
					Endpoint: "/index.php/",
					Service:  "io.quad4.web.frontend",
				},
				{
					Endpoint: "/apps/",
					Service:  "io.quad4.web.frontend",
				},
				{
					Endpoint:    "/app/list",
					Service:     "io.quad4.web.frontend",
					Unprotected: true,
				},
				{
					Endpoint: "/app/", // /app or /apps? frontend only handles /apps
					Service:  "io.quad4.web.frontend",
				},
				{
					Endpoint: "/graph/v1beta1/extensions/org.libregraph/activities",
					Service:  "io.quad4.web.activitylog",
				},
				{
					Endpoint: "/graph/v1.0/invitations",
					Service:  "io.quad4.web.invitations",
				},
				{
					Endpoint:    "/graph/v1beta1/extensions/org.libregraph/guestLinks",
					Service:     "io.quad4.web.auth-guest",
					Unprotected: true,
				},
				{
					Endpoint: "/graph/",
					Service:  "io.quad4.web.graph",
				},
				{
					Endpoint: "/api/v0/settings",
					Service:  "io.quad4.web.settings",
				},
				{
					Endpoint: "/auth-app/tokens",
					Service:  "io.quad4.web.auth-app",
				},
				{
					Endpoint:         "/wopi",
					Service:          "io.quad4.web.collaboration",
					Unprotected:      true,
					SkipXAccessToken: true,
				},
				{
					Endpoint: "/collaboration",
					Service:  "io.quad4.web.collaboration",
				},
			},
		},
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "error"
	}

	if cfg.OIDC.UserinfoCache == nil && cfg.Commons != nil && cfg.Commons.Cache != nil {
		cfg.OIDC.UserinfoCache = &config.Cache{
			Store: cfg.Commons.Cache.Store,
			Nodes: cfg.Commons.Cache.Nodes,
		}
	} else if cfg.OIDC.UserinfoCache == nil {
		cfg.OIDC.UserinfoCache = &config.Cache{}
	}

	if cfg.MachineAuthAPIKey == "" && cfg.Commons != nil && cfg.Commons.MachineAuthAPIKey != "" {
		cfg.MachineAuthAPIKey = cfg.Commons.MachineAuthAPIKey
	}

	if cfg.Reva == nil && cfg.Commons != nil {
		cfg.Reva = structs.CopyOrZeroValue(cfg.Commons.Reva)
	}

	if cfg.TransferSecret == "" && cfg.Commons != nil && cfg.Commons.TransferSecret != "" {
		cfg.TransferSecret = cfg.Commons.TransferSecret
	}

	if cfg.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}
}

// Sanitize sanitizes the configuration
func Sanitize(cfg *config.Config) {
	if cfg.Policies == nil {
		cfg.Policies = mergePolicies(DefaultPolicies(), cfg.AdditionalPolicies)
	}

	if cfg.PolicySelector == nil {
		cfg.PolicySelector = &config.PolicySelector{
			Static: &config.StaticSelectorConf{
				Policy: "default",
			},
		}
	}

	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}

	// if the CSP config file path is not set, we check if the default file exists and set it if it does
	if cfg.CSPConfigFileLocation == "" {
		defaultCSPConfigFilePath := filepath.Join(defaults.BaseDataPath(), "proxy", "csp.yaml")
		if _, err := os.Stat(defaultCSPConfigFilePath); err == nil {
			cfg.CSPConfigFileLocation = defaultCSPConfigFilePath
		}
	}
}

func mergePolicies(policies []config.Policy, additionalPolicies []config.Policy) []config.Policy {
	for _, p := range additionalPolicies {
		found := false
		for i, po := range policies {
			if po.Name == p.Name {
				po.Routes = append(po.Routes, p.Routes...)
				policies[i] = po
				found = true
				break
			}
		}
		if !found {
			policies = append(policies, p)
		}
	}
	return policies
}
