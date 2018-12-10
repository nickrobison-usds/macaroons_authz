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

build: client server endpoint

# CLI client

client: cli/build/cli

cli/build/cli:
		-mkdir -p cli/build
		-cmake -S cli -B cli/build -DOPENSSL_ROOT_DIR=/usr/local/opt/openssl
		-make -C cli/build cli

# Go Server

server: go-dep darwin/amd64

go-dep:
	dep ensure

$(PLATFORMS):
		GOOS=$(os) GOARCH=$(arch) buffalo build bin/$(NAME)_$(os)

# Javascript endpoint

endpoint: javascript/src/app.js

javascript/src/app.js:
		tsc --build javascript/tsconfig.json

.PHONY: build client clean deploy $(PLATFORMS) go-dep endpoint

