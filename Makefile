VERSION := 0.0.1
LICENSE := MIT
MAINTAINER := Nick Robison <nicholas.a.robison@omb.eop.gov>
NAME := cms_authz
PLATFORMS := darwin/amd64 linux/amd64

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

clean:
		-rm -rf cli/build
		-rm -rf bin
		-rm -r javascript/src/*.js
		-rm -rf javascript/dist

# Deploy builds

deploy: deploy-server deploy-cfssl deploy-target-service

deploy-server: linux/amd64
		packer build packer/cms_authz.json

deploy-cfssl:
		packer build packer/cfssl.json

deploy-target-service: javascript/dist/target_service.js
		packer build packer/target_service.json

.PHONY: deploy deploy-server deploy-cfssl deploy-target-service run

run:
		-cd terraform/dev; terraform apply
stop:
		-cd terraform/dev; terraform destroy

# Local client builds

build: client server endpoint

# CLI client

client: cli/build/cli

cli/build/cli:
		-mkdir -p cli/build
		-cmake -S cli -B cli/build -DOPENSSL_ROOT_DIR=/usr/local/opt/openssl
		-make -C cli/build cli

# Go Server

server: darwin/amd64

$(PLATFORMS):
		GOOS=$(os) GOARCH=$(arch) buffalo build -o bin/$(NAME)_$(os)

# Javascript endpoint

endpoint: javascript/src/app.js

javascript/dist/target_service.js:
		npm run --prefix javascript build

.PHONY: build client server clean deploy $(PLATFORMS) endpoint

