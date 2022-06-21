SHELL := bash
CURRENT_DIR = $(shell pwd)
ENV_FILE := $(CURRENT_DIR)/.env.sk

#------------------------------------------------------------------------
##help: 	                     print this help message
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/##//'

#------------------------------------------------------------------------
##test: 		     test the plugin
.PHONY:	test
test:
	@drone exec --env-file=$(ENV_FILE) --include=test

.PHONY: clean-up
#------------------------------------------------------------------------
##clean-up: 		     cleans the build artifacts
clean-up:
	rm -rf $(CURRENT_DIR)/release 

.PHONY:	build-and-push
#------------------------------------------------------------------------
#build-and-push: 		     builds and pushes the plugin image
build-and-push:	
	@drone exec --env-file=$(ENV_FILE) \
		--include=build \
		--include=publish \
		--include=publish_arm \
		--include=manifest

.PHONY:	upgrade
#------------------------------------------------------------------------
#upgrade: 		     upgrades the pipeline lib/model
upgrade:	
	curl https://raw.githubusercontent.com/drone/boilr-plugin/master/template/plugin/pipeline.go --output $(CURRENT_DIR)/plugin/pipeline.go
