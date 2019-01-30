# Macaroons Auth Demo

This project serves as a demonstration of using [Macaroons](http://macaroons.io) as the authentication mechanism for a decentralized system.

The primary context is for applications within the *Centers for Medicare and Medicaid Services*, but the majority of the concepts should be transferable to other domains as well.


## Setup

The project has a few requirements and dependencies that need to be configured correctly.
You can do this automatically by running `make setup`.
This will install all the system, javascript, and go dependencies, and initialize the terraform modules.

You can also do everything manually.

### Cloning

Since this application makes use of Go for the authentication and proxy servers, you need to clone the repo inside your `GOPATH`. This will hopefully be resolved once [this|https://github.com/nickrobison-usds/macaroons_authz/issues/1] issue has been merged.

```bash
mkdir -p $GOPATH/src/github.com/nickrobison-usds
cd $GOPATH/src/github.com/nickrobison-usds
git clone git@github.com:nickrobison-usds/macaroons_authz.git
```

We use git submodules for a number of external dependencies (to avoid requiring system installation).
You can initialize them all by running:

`git submodule update --init --recursive`

### System dependencies

We require a number of system dependencies, which are not vendored into the source tree.
The `make deps/system` command will do the installation automatically (on MacOS).

- ansible
- buffalo
- cmake
- cpprestsdk
- maven
- node
- openssl
- packer
- postgres
- terraform
- yarn

The `make` command will not install postgres by default, because the main developer prefers to use [Postgres.app](https://postgresapp.com).
A quick `brew install postgresql` should take care of that.


#### MacOS Manual installation
```bash
brew tap gobuffalo/tap
brew install ansible terraform packer cmake cpprestsdk node buffalo yarn
```

On MacOS, we cannot install Docker automatically, so you'll need to install it yourself, following the instructions [here](https://docs.docker.com/docker-for-mac/install/).


### Javascript dependencies

We use [yarn](https://yarnpkg.com) to track all of our Javascript dependencies.
The command is `make deps/js`, or the manual option:

```bash
# Install Buffalo javascript dependencies
yarn install
# Install Javascript client dependencies
cd javascript
yarn install
```

The javascript application relies on a custom fork of the `js-macaroons` dependency, which can be installed via the `npm link` command.

```bash
git clone https://github.com/nickrobison-usds/js-macaroon.git
cd js-macaroon
yarn link
```

In order for Yarn to make use of the linked package, you need to manually link the local dependency.

```bash
cd javascript
yarn link macaroon
```


### Java dependencies

Dependencies for the Java application is handled automatically by Maven.
The one caveat is that we currently require installing a custom fork of the `jmacaroons` dependency.

```bash
git clone https://github.com/nickrobison-usds/jmacaroons.git
cd jmacaroons
git checkout origin/v2-json
mvn install
```

From there, everything proceeds as normal.


### Go dependencies

Go dependencies are handled by [dep](https://golang.github.io/dep/).
They can be installed by running `make deps/go` or manually:

```bash
dep ensure
```


### Building the Components

This project features a number of distinct components (services) which serve to demonstrate a number use cases and interaction models for Macaroons.
This section briefly describes these components, along with their build instructions.


#### Authorization Server

The core project component, is the the Authorization service. This service provides the ability for assigning users to organizations, and establishing relationships between the various orgs.

The service is written in Golang, with the the [buffalo](https://gobuffalo.io/en) toolkit, which means deployment is really simple.
The `make build/server` command handles everything for you, but it currently only builds the MacOS application.

If you need to build things manually the `darwin/amd64` and `linux/amd64` make targets will build for the appropriate platform.

Also, `buffalo build` runs the default build process for the platform it's running aginst.


#### Internal Service

The *Internal Service* is an API endpoint that demonstrates an authorization flow where both the target service and the authorization service have trusted access to each other (e.g. both are on the same internal network, or managed by the same team).
This simplifies the interaction model as each service can make use of the same shared keys, which eliminates the need for public key cryptography.

This service is written in Typescript and built using [webpack](https://webpack.js.org), which compiles the source files into javascript and bundles them into a single file.

The `make build/internal` target runs the commands:

```bash
npm run --prefix javscript build
```


#### External Service

The *External Service* emulates a standalone API endpoint, which uses Macaroons to provide authentication, but without privilaged access the authorization service (e.g. a service operated by a third party, which can only communicate via the public internet).
It relies on public key cryptography to discharge third-party caveats between the services.
The emulates a fully decoupled system where the endpoint only knows a fraction of the information required to authorize a given request.


The external service is written in Java, using the [Dropwizard](https://www.dropwizard.io) framework.
It can be built automatically via the `make deploy/external-service` command, or manually with maven:

```bash
mvn package -Dmaven.javadoc.skip=true -f java/pom.xml
```

Note, Javadoc generation must be disabled under JDK 11, due to a NullPointerException that gets thrown.


#### CLI Client

The project also containts a demo command line application, written in C++, that provides some commands for interacting with the auth, internal and external services.
It simply demonstrates whether a given user is authorized to retrieve data on behalf of the organization.

The client is built with [cmake](https://cmake.org), either via the `make build/client` target, or via the manual commands.

```bash
cd cli
# We don't allow building from within the source tree
make build
cd build
cmake ..
make -j{all the cores}
```

There may be some issues with CMake finding `openssl`, mostly because `cpprestsdk` creates their own find module, which biases towards the default Homebrew location.
If that happens, you can add the `OPENSSL_ROOT_DIR` parameter to CMake.
The configure command would then become: `cmake -D OPENSSL_ROOT_DIR=/path/to/openssl ..`


### Configuration

Developer specific configuration is handled by a `.env` in the repository root.
This file is not commited to git, so it will need to be created manually before running for the first time.

Here are the necessary contents:

```bash
CFSSL_URL={http://cfssl.application.url}
```

#### Github Authentication

Application can use Github as an OAuth provider for supporting logins.
To enable, you'll need to add your credentials to the `.env` file.

```bash
CLIENT_ID={optional ClientID for Github authentication}
GITHUB_KEY={optional key for Github authentication}
GITHUB_SECRET={optional secret key for Github authentication}
```

#### Login.gov

This project supports a working, local installation of the Login.gov service.

You can clone the repo and build the their docker image, like so:

```bash
git clone https://github.com/18F/identity-idp.git
cd identity-idp
bin/setup --docker
```

You can then start the application with Login.gov enabled by adding `PROVIDER_URL={http://path.to.login.gov:local_port}` to your `.env` file.

## Run


### Running CFSSL

The Go server requires a running instance of [CFSSL](https://github.com/cloudflare/cfssl), in order to generate and manage the required public/private keys.
You can either run it via the Docker image (built with [Packer](#packer-images)), or locally.

```bash
cd cfssl
# Initialize the CA keys (only required for the first run)
cfssl genkey -initca keys/csr_ROOT_CA.json | cfssljson -bare keys/ca
# Run the CFSSL server
cfssl serve -config config/config_ca.json -ca keys/ca.pem -ca-key keys/ca-key.pem
```

### DB Seeding

The initial database can be created and populated by running the makefile commands.

The `database.yml` contains the necessary configuration options, which can be set to the appropriate values for the developer's environment.
At a minimum, you'll need to modify the `user` parameter to match the username used by Postgres (most likely the user's local account name).

```yaml
development:
  dialect: postgres
  database: macaroons_authz_development
  user: {DATABASE USERNAME}
  password: {DATABASE PASSWORD (if required)} 
  host: 127.0.0.1
  pool: 5
```

This is done manually, to avoid destroying any existing data.

```bash
make build/database
make build/seed
```

Note: There's currently a bug in the implementation where the seeding script does not properly remove the `root_keys` table.
This means you'll need to manually remove it each time you want to re-run the seeding process.

You can also manually run the seeding by the `buffalo task db:seed` command.


### Development modes

You can run both the Go server and the Javascript client in dev mode through the following commands:

`PORT={go server port} buffalo dev` from the repo root will start the Go server in dev mode, which means it will automatically reload when any of its source files change.

`npm run -c javascript watch-debug` performs the same actions, but for the Javascript client.


## Deploy

Each service can be run directly on the developer's machine, or within an isolated Docker environment, provided by [Terraform](https://www.terraform.io). 
The Docker environment is the recommended way of standing everything up.
We currently support an `sbx` (sandbox) environment which doesn't persist any data to disk.


### Build Packer images

The Docker images are built using [Packer](https://www.packer.io) with the setup scripts making use of [Ansible](https://www.ansible.com).
Each service can be built by calling `packer build` on each file in the `packer/` directory.

Of course, the Makefile will handle it all for you, via the `build/deploy` target, which rebuilds all of the services.
It also handles generating the required binaries, which are then copied into the Docker images.


## Launching Docker/Terraform

Running everything can be done via the `run` target in the Makefile, likewise `stop` shuts everything down and removes the temporary data.

