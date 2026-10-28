# IDP

This service provides a builtin minimal OpenID Connect provider based on [LibreGraph Connect (lico)](https://github.com/libregraph/lico) for Quad4 Cloud.

It is mainly targeted at smaller installations. For larger setups it is recommended to replace IDP with an external OpenID Connect Provider.

By default, it is configured to use the Quad4 Cloud IDM service as its LDAP backend for looking up and authenticating users. Other backends like an external LDAP server can be configured via a set of [enviroment variables](https://github.com/Quad4-Software/cloud/docs/dev/server/services/idp/environment-variables).

Note that translations provided by the IDP service are not maintained via Quad4 Cloud but part of the embedded  [LibreGraph Connect Identifier](https://github.com/libregraph/lico/tree/master/identifier) package.

## Configuration

### Custom Clients

By default the `idp` service generates a OIDC client configuration suitable for
using Quad4 Cloud with the standard client applications (Web, Desktop, iOS and
Android). If you need to configure additional client it is possible to inject a
custom configuration via `yaml`. This can be done by adding a section `clients`
to the `idp` section of the main configuration file (`quad4.yaml`). This section
needs to contain configuration for all clients (including the standard clients).

For example if you want to add a (public) client for use with the oidc-agent you would
need to add this snippet to the `idp` section in `quad4.yaml`.

```yaml
clients:
- id: web
  name: Quad4 Cloud Web App
  trusted: true
  secret: ""
  redirect_uris:
  - https://quad4.k8s:9200/
  - https://quad4.k8s:9200/oidc-callback.html
  - https://quad4.k8s:9200/oidc-silent-redirect.html
  post_logout_redirect_uris: []
  origins:
  - https://quad4.k8s:9200
  application_type: ""
- id: Quad4 CloudDesktop
  name: Quad4 Cloud Desktop Client
  trusted: false
  secret: ""
  redirect_uris:
  - http://127.0.0.1
  - http://localhost
  post_logout_redirect_uris: []
  origins: []
  application_type: native
- id: Quad4 CloudAndroid
  name: Quad4 Cloud Android App
  trusted: false
  secret: ""
  redirect_uris:
  - oc://android.quad4.io
  post_logout_redirect_uris:
  - oc://android.quad4.io
  origins: []
  application_type: native
- id: Quad4 CloudIOS
  name: Quad4 Cloud iOS App
  trusted: false
  secret: ""
  redirect_uris:
  - oc://ios.quad4.io
  post_logout_redirect_uris:
  - oc://ios.quad4.io
  origins: []
  application_type: native
- id: oidc-agent
  name: OIDC Agent
  trusted: false
  secret: ""
  redirect_uris:
  - http://127.0.0.1
  - http://localhost
  post_logout_redirect_uris: []
  origins: []
  application_type: native
```



