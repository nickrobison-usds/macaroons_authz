/* ----- Setup the Login.gov resource ----- */
module "db" {
  source = "../modules/db"

  name = "cmsauthz"
  host_path = "/Users/usds/Development/identity-idp/db_data"
  use_local = false
  db_name = "cms_authz"
}
/*
module "idp" {
  source = "../modules/idp"

  name = "idp"
  db_hostname = "${module.db.hostname}"
  db_network = "${module.db.network}"
  public_network = "${docker_network.public.name}"
  host_path = "/Users/usds/Development/identity-idp"
}
*/

module "cfssl" {
  source = "../modules/cfssl"

  name = "cfssl"
  public_network = "${docker_network.public.name}"
}

/* ----- Setup the AuthZ resource ----- */

// Create the database
/*
provider "postgresql" {
  alias = "pg1"
  host = "${module.db.hostname}"
  username = "postgres"
  password = "postgres"
  sslmode = "disable"
}

resource "postgresql_database" "authz_db" {
  provider = "postgresql.pg1"
  name = "cmsauthz-db"
}
*/

resource "docker_container" "authz" {
  name = "authz_server"
  image = "nickrobison.com/cms_authz:latest"
  hostname = "server"
  ports {
    internal = 8080
    external = 8080
  }
  env = [
    "DATABASE_URL=postgres://postgres@${module.db.hostname}:5432/cms_authz?sslmode=disable",
    "CFSSL_URL=http://${module.cfssl.hostname}:8888",
    "PROVIDER_URL=http://localhost:3000",
    "PORT=8080",
    "SEED=true",
    "GO_ENV=production"
  ]
  networks_advanced {
    name = "${docker_network.public.name}"
  }
  networks_advanced {
    name = "${module.db.network}"
  }
}

/* ----- Setup the Target Service resource ----- */

resource "docker_container" "target-service" {
  name = "target-service"
  image = "nickrobison.com/target_service:latest"
  hostname = "target"
  ports {
    internal = 3002
    external = 3002
  }
  env = [
    "DATABASE_URL=postgres://postgres@${module.db.hostname}:5432/cms_authz?sslmode=disable",
  ]
  networks_advanced {
    name = "${docker_network.public.name}"
  }
  networks_advanced {
    name = "${module.db.network}"
  }
}

/* ----- Public networking ----- */

resource "docker_network" "public" {
  name = "cms_authz-public"
}
