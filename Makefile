# Copyright (c) Mainflux
# SPDX-License-Identifier: Apache-2.0

# Adapted for Orb project, modifications licensed under MPL v. 2.0:
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/. */

MF_DOCKER_IMAGE_NAME_PREFIX ?= orb
DOCKERHUB_REPO = ns1labs
BUILD_DIR = build
SERVICES = fleet policies sinks prom-sink agent
DOCKERS = $(addprefix docker_,$(SERVICES))
DOCKERS_DEV = $(addprefix docker_dev_,$(SERVICES))
CGO_ENABLED ?= 0
GOARCH ?= amd64

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=mod -ldflags "-s -w" -o ${BUILD_DIR}/$(MF_DOCKER_IMAGE_NAME_PREFIX)-$(1) cmd/$(1)/main.go
endef

define compile_service_linux
	$(eval svc=$(subst docker_dev_,,$(1)))
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=mod -ldflags "-s -w" -o ${BUILD_DIR}/$(MF_DOCKER_IMAGE_NAME_PREFIX)-$(svc) cmd/$(svc)/main.go
endef

define make_docker
	$(eval svc=$(subst docker_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--tag=$(DOCKERHUB_REPO)/$(MF_DOCKER_IMAGE_NAME_PREFIX)-$(svc) \
		-f docker/Dockerfile .
endef

define make_docker_dev
	$(eval svc=$(subst docker_dev_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--tag=$(DOCKERHUB_REPO)/$(MF_DOCKER_IMAGE_NAME_PREFIX)-$(svc) \
		-f docker/Dockerfile.dev ./build
endef

all: $(SERVICES)

.PHONY: all $(SERVICES) dockers dockers_dev latest release

clean:
	rm -rf ${BUILD_DIR}

cleandocker:
	# Stops containers and removes containers, networks, volumes, and images created by up
#	docker-compose -f docker/docker-compose.yml down --rmi all -v --remove-orphans
	docker-compose -f docker/docker-compose.yml down -v --remove-orphans

ifdef pv
	# Remove unused volumes
	docker volume ls -f name=$(MF_DOCKER_IMAGE_NAME_PREFIX) -f dangling=true -q | xargs -r docker volume rm
endif

install:
	cp ${BUILD_DIR}/* $(GOBIN)

test:
	go test -mod=mod -race -count 1 -tags test $(shell go list ./... | grep -v 'cmd')

proto:
	protoc --gofast_out=plugins=grpc:. *.proto

$(SERVICES):
	$(call compile_service,$(@))

$(DOCKERS):
	$(call make_docker,$(@),$(GOARCH))

$(DOCKERS_DEV):
	$(call compile_service_linux,$(@))
	$(call make_docker_dev,$(@))

dockers: $(DOCKERS)
dockers_dev: $(DOCKERS_DEV)

define docker_push
	for svc in $(SERVICES); do \
		docker push $(DOCKERHUB_REPO)/$(MF_DOCKER_IMAGE_NAME_PREFIX)-$$svc:$(1); \
	done
endef

run:
	docker-compose -f docker/docker-compose.yml up
