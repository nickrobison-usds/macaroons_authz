VERSION := 0.0.1
LICENSE := MIT
MAINTAINER := Nick Robison <nicholas.a.robison@omb.eop.gov>
NAME := Macaroons Authorization Demo
PLATFORMS := darwin/amd64 linux/amd64

# Check for required packages
UNAME := $(shell uname)
PKGS := cmake dep yarn buffalo cpprestsdk cfssl

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

clean:
		-rm -rf cli/build
		-rm -rf bin
		-rm -r javascript/src/*.js
		-rm -rf javascript/dist

#
# Setup repository for the first time
#
setup: deps/js deps/go

# Install system dependencies via the OS package manager.
deps/system:
	@echo "Installing required dependencies: $(PKGS)"
ifeq ($(UNAME),Linux)
	@echo "Linux deps aren't installed automatically"
endif
ifeq ($(UNAME),Darwin)
	brew tap gobuffalo/tap
	brew install $(PKGS)
endif

# Install the required Javascript dependencies via Yarn
deps/js: deps/js/client deps/js/server

# Install all the required Javascript dependencies for the client application
deps/js/client: deps/system
	yarn --cwd javascript install

# Install all the required Javascript dependencies for the Buffalo server
deps/js/server: deps/system
	yarn install

# Install required Go dependencies
deps/go: deps/system
	dep ensure

.PHONY: setup deps/system deps/js deps/js/client deps/js/server deps/go

#
# Local application builds
#

build: build/client build/server build/endpoint

# CLI client

build/client: cli/build/cli

cli/build/cli:
		-mkdir -p cli/build
		-cmake -S cli -B cli/build -DOPENSSL_ROOT_DIR=/usr/local/opt/openssl
		-make -C cli/build cli

# Go Server

build/server: darwin/amd64

$(PLATFORMS):
		GOOS=$(os) GOARCH=$(arch) buffalo build -o bin/macaroons_authz_$(os)

# Javascript endpoint

build/endpoint: javascript/dist/target_service.js

javascript/dist/target_service.js:
		npm run --prefix javascript build

.PHONY: build client server clean deploy $(PLATFORMS) endpoint


# Deploy builds

deploy: deploy-server deploy-cfssl deploy-target-service

deploy-server: linux/amd64
		packer build packer/macaroons_authz.json

deploy-cfssl:
		packer build packer/cfssl.json

deploy-target-service: javascript/dist/target_service.js
		packer build packer/target_service.json

.PHONY: deploy deploy-server deploy-cfssl deploy-target-service run

run:
		-cd terraform/sbx; terraform apply

stop:
		-cd terraform/sbx; terraform destroy


