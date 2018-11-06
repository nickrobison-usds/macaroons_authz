/* ----- Setup the Login.gov resource ----- */
module "db" {
  source = "../modules/db"

  name = "cmsauthz"
  host_path = "/Users/usds/Development/identity-idp/db_data"
  use_local = false
}

module "idp" {
  source = "../modules/idp"

  name = "idp"
  db_hostname = "${module.db.hostname}"
  db_network = "${module.db.network}"
  public_network = "${docker_network.public.name}"
  host_path = "/Users/usds/Development/identity-idp"
}

resource "docker_network" "public" {
  name = "cms_authz-public"
}
