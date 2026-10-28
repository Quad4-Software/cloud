# Changelog

All notable changes to Quad4 Cloud are documented in this file.

## Unreleased

Quad4 Cloud is a fork of OpenCloud, which is itself a fork of ownCloud
Infinite Scale. The fork point is upstream commit 5d745311.

### Changed

- Rebranded the product, web interface, login screen, emails, and docs
  to Quad4 Cloud.
- The Go module is now `github.com/Quad4-Software/cloud` and the binary
  is `quad4`.
- Service registry names moved from `eu.opencloud.*` to `io.quad4.*`.
- Default config path is `/etc/quad4`, default data path is
  `/var/lib/quad4`, home-dir based paths use `~/.quad4`.
- Container image is built on Chainguard wolfi-base, runs as the
  unprivileged `nonroot` user, and is published to
  `ghcr.io/quad4-software/cloud` with zstd-compressed layers and cosign
  keyless signatures.
- CI is GitHub Actions; the upstream Woodpecker configuration was
  removed.

### Compatibility notes

- Environment variable prefix remains `OC_` and `OCIS_` for now.
- Deployments migrating from OpenCloud must move config and data
  directories or set `OC_CONFIG_DIR` and `OC_BASE_DATA_PATH`.
