VERSION := 0.0.1
LICENSE := MIT
MAINTAINER := Nick Robison <nicholas.a.robison@omb.eop.gov>
NAME := Macaroons Authorization Demo
PLATFORMS := darwin/amd64 linux/amd64

# Check for required packages
UNAME := $(shell uname)
PKGS := cmake dep yarn buffalo cpprestsdk cfssl maven ansible packer terraform

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

#
# Setup repository for the first time
#
setup: deps/js deps/go deps/python

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
	-cd terraform/sbx; terraform init

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

deps/python: deps/system
	pip install sphinx sphinx-rtd-theme

.PHONY: setup deps/system deps/js deps/js/client deps/js/server deps/go deps/python

#
# Local application builds
#

build: build/client build/server build/endpoint build/seed

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
		GOOS=$(os) GOARCH=$(arch) go build -o bin/proxy_server_$(os) proxy/main.go


# Javascript endpoint

build/internal: javascript/dist/target_service.js

javascript/dist/target_service.js:
		npm run --prefix javascript build

# Java client (not currently working)
java/target/javaservice-%.jar:
	mvn package -Dmaven.javadoc.skip=true -f java/pom.xml

build/dependencies:
	go get github.com/gobuffalo/buffalo-pop

build/database: build/dependencies
	buffalo db create -a
	buffalo db migrate

build/seed:
	buffalo task db:seed

.PHONY: build build/client build/server $(PLATFORMS) build/internal build/seed build/dependencies build/external


# Deploy builds

deploy: deploy/server deploy/cfssl deploy/internal-service deploy/external-service

deploy/server: linux/amd64
		packer build packer/macaroons_authz.json

deploy/cfssl:
		packer build packer/cfssl.json

deploy/internal-service: javascript/dist/target_service.js
		packer build packer/internal_service.json

deploy/external-service: java/target/javaservice-%.jar
		packer build packer/external_service.json

deploy/proxy: linux/amd64
		packer build packer/proxy.json

run:
		-cd terraform/sbx; terraform apply

stop:
		-cd terraform/sbx; terraform destroy

.PHONY: deploy deploy/server deploy/cfssl deploy/internal-service deploy/external-service deploy/proxy run stop

# Documentation

docs:
	sphinx-build -b html docs/ docs/_build

.PHONY: docs

clean:
		-rm -rf cli/build
		-rm -rf bin
		-rm -r javascript/src/*.js
		-rm -rf javascript/dist
		-rm -rf java/target
