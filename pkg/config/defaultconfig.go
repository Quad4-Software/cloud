package config

import (
	"github.com/Quad4-Software/cloud/pkg/shared"
	activitylog "github.com/Quad4-Software/cloud/services/activitylog/pkg/config/defaults"
	antivirus "github.com/Quad4-Software/cloud/services/antivirus/pkg/config/defaults"
	appProvider "github.com/Quad4-Software/cloud/services/app-provider/pkg/config/defaults"
	appRegistry "github.com/Quad4-Software/cloud/services/app-registry/pkg/config/defaults"
	audit "github.com/Quad4-Software/cloud/services/audit/pkg/config/defaults"
	authapp "github.com/Quad4-Software/cloud/services/auth-app/pkg/config/defaults"
	authbasic "github.com/Quad4-Software/cloud/services/auth-basic/pkg/config/defaults"
	authbearer "github.com/Quad4-Software/cloud/services/auth-bearer/pkg/config/defaults"
	authguest "github.com/Quad4-Software/cloud/services/auth-guest/pkg/config/defaults"
	authmachine "github.com/Quad4-Software/cloud/services/auth-machine/pkg/config/defaults"
	authservice "github.com/Quad4-Software/cloud/services/auth-service/pkg/config/defaults"
	clientlog "github.com/Quad4-Software/cloud/services/clientlog/pkg/config/defaults"
	collaboration "github.com/Quad4-Software/cloud/services/collaboration/pkg/config/defaults"
	eventhistory "github.com/Quad4-Software/cloud/services/eventhistory/pkg/config/defaults"
	frontend "github.com/Quad4-Software/cloud/services/frontend/pkg/config/defaults"
	gateway "github.com/Quad4-Software/cloud/services/gateway/pkg/config/defaults"
	graph "github.com/Quad4-Software/cloud/services/graph/pkg/config/defaults"
	groups "github.com/Quad4-Software/cloud/services/groups/pkg/config/defaults"
	idm "github.com/Quad4-Software/cloud/services/idm/pkg/config/defaults"
	idp "github.com/Quad4-Software/cloud/services/idp/pkg/config/defaults"
	invitations "github.com/Quad4-Software/cloud/services/invitations/pkg/config/defaults"
	nats "github.com/Quad4-Software/cloud/services/nats/pkg/config/defaults"
	notifications "github.com/Quad4-Software/cloud/services/notifications/pkg/config/defaults"
	ocm "github.com/Quad4-Software/cloud/services/ocm/pkg/config/defaults"
	ocs "github.com/Quad4-Software/cloud/services/ocs/pkg/config/defaults"
	policies "github.com/Quad4-Software/cloud/services/policies/pkg/config/defaults"
	postprocessing "github.com/Quad4-Software/cloud/services/postprocessing/pkg/config/defaults"
	proxy "github.com/Quad4-Software/cloud/services/proxy/pkg/config/defaults"
	search "github.com/Quad4-Software/cloud/services/search/pkg/config/defaults"
	settings "github.com/Quad4-Software/cloud/services/settings/pkg/config/defaults"
	sharing "github.com/Quad4-Software/cloud/services/sharing/pkg/config/defaults"
	sse "github.com/Quad4-Software/cloud/services/sse/pkg/config/defaults"
	storagepublic "github.com/Quad4-Software/cloud/services/storage-publiclink/pkg/config/defaults"
	storageshares "github.com/Quad4-Software/cloud/services/storage-shares/pkg/config/defaults"
	storageSystem "github.com/Quad4-Software/cloud/services/storage-system/pkg/config/defaults"
	storageusers "github.com/Quad4-Software/cloud/services/storage-users/pkg/config/defaults"
	thumbnails "github.com/Quad4-Software/cloud/services/thumbnails/pkg/config/defaults"
	userlog "github.com/Quad4-Software/cloud/services/userlog/pkg/config/defaults"
	users "github.com/Quad4-Software/cloud/services/users/pkg/config/defaults"
	web "github.com/Quad4-Software/cloud/services/web/pkg/config/defaults"
	webdav "github.com/Quad4-Software/cloud/services/webdav/pkg/config/defaults"
	webfinger "github.com/Quad4-Software/cloud/services/webfinger/pkg/config/defaults"
)

func DefaultConfig() *Config {
	return &Config{
		Quad4CloudURL: "https://localhost:9200",
		Runtime: Runtime{
			Port:          "9250",
			Host:          "localhost",
			ShutdownOrder: []string{"proxy"},
		},
		Reva: &shared.Reva{
			Address: "io.quad4.api.gateway",
		},

		Activitylog:       activitylog.DefaultConfig(),
		Antivirus:         antivirus.DefaultConfig(),
		AppProvider:       appProvider.DefaultConfig(),
		AppRegistry:       appRegistry.DefaultConfig(),
		Audit:             audit.DefaultConfig(),
		AuthApp:           authapp.DefaultConfig(),
		AuthBasic:         authbasic.DefaultConfig(),
		AuthBearer:        authbearer.DefaultConfig(),
		AuthMachine:       authmachine.DefaultConfig(),
		AuthService:       authservice.DefaultConfig(),
		Clientlog:         clientlog.DefaultConfig(),
		Collaboration:     collaboration.DefaultConfig(),
		EventHistory:      eventhistory.DefaultConfig(),
		Frontend:          frontend.DefaultConfig(),
		Gateway:           gateway.DefaultConfig(),
		Graph:             graph.DefaultConfig(),
		Groups:            groups.DefaultConfig(),
		AuthGuest:         authguest.DefaultConfig(),
		IDM:               idm.DefaultConfig(),
		IDP:               idp.DefaultConfig(),
		Invitations:       invitations.DefaultConfig(),
		Nats:              nats.DefaultConfig(),
		Notifications:     notifications.DefaultConfig(),
		OCM:               ocm.DefaultConfig(),
		OCS:               ocs.DefaultConfig(),
		Postprocessing:    postprocessing.DefaultConfig(),
		Policies:          policies.DefaultConfig(),
		Proxy:             proxy.DefaultConfig(),
		Search:            search.DefaultConfig(),
		Settings:          settings.DefaultConfig(),
		Sharing:           sharing.DefaultConfig(),
		SSE:               sse.DefaultConfig(),
		StoragePublicLink: storagepublic.DefaultConfig(),
		StorageShares:     storageshares.DefaultConfig(),
		StorageSystem:     storageSystem.DefaultConfig(),
		StorageUsers:      storageusers.DefaultConfig(),
		Thumbnails:        thumbnails.DefaultConfig(),
		Userlog:           userlog.DefaultConfig(),
		Users:             users.DefaultConfig(),
		Web:               web.DefaultConfig(),
		WebDAV:            webdav.DefaultConfig(),
		Webfinger:         webfinger.DefaultConfig(),
	}
}
