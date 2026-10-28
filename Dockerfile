# syntax=docker/dockerfile:1
# Quad4 Cloud container image.
#
# Build locally with podman:
#   podman build -t quad4-cloud:dev .
#
# Multi-stage build on Chainguard wolfi-base. The final image runs as the
# unprivileged nonroot user (uid 65532).

FROM cgr.dev/chainguard/wolfi-base:latest AS nodebuild

RUN apk add --no-cache bash curl git make nodejs-24 pnpm-11

WORKDIR /src
COPY ./ /src/

RUN make node-generate-prod

FROM cgr.dev/chainguard/wolfi-base:latest AS build

RUN apk add --no-cache bash build-base git go-1.26 libvips-dev make pkgconf

WORKDIR /src
COPY --from=nodebuild /src /src

ARG VERSION=dev
ARG EDITION=dev
RUN make -C quad4 go-generate release-linux-docker-$(go env GOARCH) ENABLE_VIPS=true VERSION=${VERSION} EDITION=${EDITION}

FROM cgr.dev/chainguard/wolfi-base:latest

RUN apk add --no-cache attr ca-certificates curl mailcap tree libvips && \
	echo 'hosts: files dns' >| /etc/nsswitch.conf && \
	mkdir -p /etc/quad4 /var/lib/quad4 && \
	chown -R nonroot:nonroot /etc/quad4 /var/lib/quad4

LABEL org.opencontainers.image.title="Quad4 Cloud" \
	org.opencontainers.image.vendor="Quad4" \
	org.opencontainers.image.authors="Quad4" \
	org.opencontainers.image.description="Quad4 Cloud is a modern file-sync and share platform" \
	org.opencontainers.image.licenses="Apache-2.0" \
	org.opencontainers.image.documentation="https://github.com/Quad4-Software/cloud" \
	org.opencontainers.image.source="https://github.com/Quad4-Software/cloud" \
	org.opencontainers.image.url="https://quad4.io"

USER nonroot

VOLUME ["/etc/quad4", "/var/lib/quad4"]

EXPOSE 9200

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
CMD ["server"]

COPY --from=build /src/quad4/dist/binaries/quad4-linux-* /usr/bin/quad4
COPY --chmod=0755 entrypoint.sh /usr/bin/entrypoint.sh
