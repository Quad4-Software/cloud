![Quad4 Cloud](quad4_logo.png)

# Quad4 Cloud

Quad4 Cloud is a file sync and share platform. It provides file storage,
sharing, collaboration, and user management over a web interface, WebDAV,
and the LibreGraph API. It is a fork of OpenCloud.

## Run with containers

Multi-arch images (linux/amd64 and linux/arm64) are published to
`ghcr.io/quad4-software/cloud` with zstd-compressed layers and cosign
keyless signatures. The image runs as an unprivileged user.

```console
podman pull ghcr.io/quad4-software/cloud:latest
podman run -d --name quad4-cloud \
  -p 9200:9200 \
  -v quad4-config:/etc/quad4 \
  -v quad4-data:/var/lib/quad4 \
  ghcr.io/quad4-software/cloud:latest
```

On first start the server initializes its configuration and data, then
serves the web interface at `https://localhost:9200`. For local testing,
`OC_INSECURE=true` enables the built-in self-signed certificate:

```console
podman run -d --name quad4-cloud \
  -p 9200:9200 \
  -e OC_INSECURE=true \
  -e OC_URL=https://localhost:9200 \
  -e IDM_ADMIN_PASSWORD=admin \
  -v quad4-config:/etc/quad4 \
  -v quad4-data:/var/lib/quad4 \
  ghcr.io/quad4-software/cloud:latest
```

Log in as `admin` with the configured `IDM_ADMIN_PASSWORD`.

Verify the image signature:

```console
cosign verify ghcr.io/quad4-software/cloud:latest \
  --certificate-identity-regexp 'https://github.com/Quad4-Software/cloud' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com
```

## Configuration

Configuration lives in `/etc/quad4`, data in `/var/lib/quad4`. All
settings can also be set via environment variables; see the `env` fields
in each service's `pkg/config/config.go`.

### OpenID Connect

Quad4 Cloud ships an embedded identity provider. To use an external
OIDC provider such as Keycloak or Authentik instead:

```console
-e PROXY_OIDC_ISSUER=https://idp.example.com/realms/quad4 \
-e PROXY_OIDC_REWRITE_WELLKNOWN=true \
-e WEB_OIDC_CLIENT_ID=quad4-web \
-e PROXY_AUTOPROVISION_ACCOUNTS=true \
-e PROXY_USER_OIDC_CLAIM=preferred_username \
-e PROXY_AUTOPROVISION_CLAIM_USERNAME=preferred_username \
-e PROXY_AUTOPROVISION_CLAIM_EMAIL=email \
-e PROXY_AUTOPROVISION_CLAIM_DISPLAYNAME=name \
-e PROXY_AUTOPROVISION_CLAIM_GROUPS=groups
```

When `PROXY_OIDC_ISSUER` is set, requests to `/.well-known/openid-configuration`
are rewritten to the external provider and the embedded IDP stays unused.
Auto-provisioned users are created on first login.

### S3 storage backend

By default files are stored on the local filesystem (posix driver). To
store blobs in S3 while keeping metadata on the filesystem, use the
`decomposeds3` driver:

```console
-e STORAGE_USERS_DRIVER=decomposeds3 \
-e STORAGE_USERS_DECOMPOSEDS3_ENDPOINT=https://s3.example.com \
-e STORAGE_USERS_DECOMPOSEDS3_REGION=us-east-1 \
-e STORAGE_USERS_DECOMPOSEDS3_BUCKET=quad4 \
-e STORAGE_USERS_DECOMPOSEDS3_ACCESS_KEY=... \
-e STORAGE_USERS_DECOMPOSEDS3_SECRET_KEY=...
```

Works with AWS S3 and S3-compatible stores such as MinIO.

## Build from source

Requires Go, Node.js, and pnpm:

```console
make generate
make -C quad4 build
quad4/bin/quad4 init && quad4/bin/quad4 server
```

The container image can be built locally with podman:

```console
podman build -t quad4-cloud:dev .
```

## Security

Report security issues to [security@quad4.io](mailto:security@quad4.io).

## License

Apache 2.0. See [LICENSE](LICENSE). Quad4 Cloud is a fork of OpenCloud,
which is itself a fork of ownCloud Infinite Scale.
