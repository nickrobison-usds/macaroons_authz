# CMS Auth Demo

## Setup

The project has a few requirements and dependencies that need to be configured correctly.
You can do this automatically by running `make setup`.
This will install all the system, javascript, and go dependencies, and initialize the terraform modules.

You can also do everything manually.

### Cloning

We use git submodules for a number of external dependencies (to avoid requiring system installation).
You can initialize them all by running

`git submodule init --update --recursive`.

Or, the [command line client](#cli-client) will handle it automatically. 

### System dependencies

We require a number of system dependencies, which are not vendored into the source tree.
The `make deps/system` command will do the installation automatically (on MacOS).

- ansible
- terraform
- packer
- cmake
- cpprestsdk
- buffalo
- node
- yarn
- openssl

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

### Go dependencies

Go dependencies are handled by [dep](https://golang.github.io/dep/).
They can be installed by running `make deps/go` or manually:

```bash
dep ensure
```

### Building

#### Go server

The main go server can be 


#### CLI Client

The demo client is a C++ application that is built with [cmake](https://cmake.org)

You can build it via `make client` or via the manual commands.

```bash
cd cli
# We don't allow building from within the source tree
make build
cd build
cmake ..
make -j{all the cores}
```





### Configuration

`.env` file in the root directory

### Login.gov

This project supports a working, local installation of the Login.gov service.

You can clone the repo and build the docker images, like so:

```bash
git clone https://github.com/18F/identity-idp.git
cd identity-idp
bin/setup --docker
```

## Build

`make build`


### Packer images

Must be built from the root directory!

## Deploy

`make deploy`
