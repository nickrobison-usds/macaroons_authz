# CMS Auth Demo

## Setup

### System dependencies

```bash
brew tap gobuffalo/tap
brew install ansible terraform packer cmake cpprestsdk node buffalo yarn
```

Install docker


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
