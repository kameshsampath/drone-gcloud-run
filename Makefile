IMAGE?=kameshsampath/drone-gcloud-run
TAG?=latest
SHELL := bash
CURRENT_DIR = $(shell pwd)
ENV_FILE := $(CURRENT_DIR)/.env.sk
BUILDER=buildx-multi-arch
DOCKER_FILE=$(CURRENT_DIR)/docker/Dockerfile

bin:	## Build binaries
	goreleaser build --snapshot --rm-dist --single-target --debug

bin-all:	## Build binaries for all targetted architectures
	goreleaser build --snapshot --rm-dist

tidy:	## Runs go mod tidy
	go mod tidy
	
vendor:	## Vendoring
	go mod vendor

lint:	## Run lint on the project
	golangci-lint run
		    
.PHONY:	test	# test the plugin
test:
	@drone exec --env-file=$(ENV_FILE) --include=test

clean: #cleans the build artifacts
	rm -rf $(CURRENT_DIR)/dist 

prepare-buildx: ## Create buildx builder for multi-arch build, if not exists
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host

build-plugin: ## Build plugin image locally
	docker build --tag=$(IMAGE):$(TAG) -f $(DOCKER_FILE) .
	docker tag $(IMAGE):$(TAG) $(IMAGE):$(TAG)

push-plugin: prepare-buildx ## Build & Upload extension image to hub. Do not push if tag already exists: TAG=$(svu c) make push-extension
	docker pull $(IMAGE):$(shell svu c --strip-prefix) && echo "Failure: Tag already exists" || docker buildx build --push --builder=$(BUILDER) --platform=linux/amd64,linux/arm64 --build-arg TAG=$(shell svu c --strip-prefix) --tag=$(IMAGE):$(shell svu c --strip-prefix) -f $(DOCKER_FILE) .

release:	
	git tag "$(shell svu next --strip-prefix)"
	git push --tags
	goreleaser --rm-dist

.PHONY:	upgrade
upgrade:	#upgrades the pipeline lib/model
	curl https://raw.githubusercontent.com/drone/boilr-plugin/master/template/plugin/pipeline.go --output $(CURRENT_DIR)/plugin/pipeline.go

help: ## Show this help
	@echo Please specify a build target. The choices are:
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(INFO_COLOR)%-30s$(NO_COLOR) %s\n", $$1, $$2}'

.PHONY: bin clean extension push-plugin help	tidy	test	vendor	lint	clean