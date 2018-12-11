/* ----- Setup the Login.gov resource ----- */
module "db" {
  source = "../modules/db"

  name = "cmsauthz"
  host_path = "/Users/usds/Development/identity-idp/db_data"
  use_local = false
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

/* ----- Setup the AuthZ resources ----- */

resource "docker_container" "authz" {
  name = "authz_server"
  image = "nickrobison.com/cms_authz:latest"
  hostname = "server"
  ports {
    internal = 8080
    external = 8080
  }
  env = [
    "DATABASE_URL=postgres://postgres@${module.db.hostname}?sslmode=disable",
    "CFSSL_URL=http://${module.cfssl.hostname}",
    "PROVIDER_URL=http://localhost:3000",
    "PORT=8080"
  ]
  networks_advanced {
    name = "${docker_network.public.name}"
  }
  networks_advanced {
    name = "${module.db.network}"
  }
}

/* ----- Setup the AuthZ resources ----- */

resource "docker_container" "target-service" {
  name = "target-service"
  image = "nickrobison.com/target_service:latest"
  hostname = "target"
  ports {
    internal = 3002
    external = 3002
  }
  networks_advanced {
    name = "${docker_network.public.name}"
  }
}

resource "docker_network" "public" {
  name = "cms_authz-public"
}
