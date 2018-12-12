/* DB setup */

locals = {
  network_name = "${var.name}-db-private"
  hostname = "${var.name}-db"
  /* This is how we can get an empty name for the local db network */
  local_network = ""
}

resource "docker_container" "postgres" {
  name = "${var.name}-db"
  count = "${1 - var.use_local}"
  hostname = "${local.hostname}"
  image = "${docker_image.postgres.latest}"
  networks_advanced {
    name = "${docker_network.db_private.name}"
  }
  volumes {
    host_path = "${var.host_path}"
    container_path ="/var/lib/postgresql/data"
  }
  env = [
    "POSTGRES_DB=${var.db_name}"
  ]
}

resource "docker_network" "db_private" {
  count = "${1 - var.use_local}"
  name = "${local.network_name}"
  internal = true
}

resource "docker_image" "postgres" {
  name = "postgres"
}
