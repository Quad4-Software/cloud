SHELL := bash
NAME := quad4
TAGS := disable_crypt
GOARCH := $(shell go env GOARCH)
CONFIG_DOCS_BASE_PATH := ../docs
ifdef ENABLE_VIPS
	TAGS := ${TAGS},enable_vips
endif

ifneq (, $(shell command -v go 2> /dev/null)) # suppress `command not found warnings` for non go targets in CI
include ../.bingo/Variables.mk
endif
include ../.make/default.mk
include ../.make/go.mk
include ../.make/release.mk
include ../.make/docs.mk

.PHONY: dev-docker
dev-docker:
	docker build -f docker/Dockerfile.multiarch -t quad4-software/quad4:dev ..

.PHONY: dev-docker-amd64 dev-docker-arm64
dev-docker-amd64 dev-docker-arm64: dev-docker-%:
	docker buildx build --platform linux/$* --load -f docker/Dockerfile.multiarch -t quad4-software/quad4:dev-$* ..

dev-docker-with-workspace:
	docker build -f docker/Dockerfile.multiarch -t quad4-software/quad4:dev --build-arg SRCDIR=quad4 ../..

.PHONY: dev-docker-multiarch
dev-docker-multiarch:
	@echo "+-------------------------------------------------------------------------+"
	@echo "| Are you sure you have run make node-generate-prod in the repository root? |"
	@echo "+-------------------------------------------------------------------------+"
	sleep 10
	docker buildx rm quad4builder || true
	docker buildx create --platform linux/arm64,linux/amd64 --name quad4builder
	docker buildx use quad4builder
	docker buildx build --platform linux/arm64,linux/amd64 --output type=docker --file docker/Dockerfile.multiarch --tag quad4-software/quad4:dev-multiarch ..
	docker buildx rm quad4builder

.PHONY: debug-docker
debug-docker:
	$(MAKE) --no-print-directory debug-linux-docker-$(GOARCH)
	docker build -f docker/Dockerfile.linux.debug.$(GOARCH) -t quad4-software/quad4:debug .
