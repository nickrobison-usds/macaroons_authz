resource "docker_container" "idp" {
  name = "idp_web"
  image = "nickrobison.com/idp-web:latest"
  hostname = "${var.name}"
  ports {
    internal = 3000
    external = 3000
  }
  env = [
    "SEED=${var.seed}",
    "REDIS_URL=redis://redis",
    "DATABASE_URL=postgres://postgres@${var.db_hostname}",
    "DOCKER_DB_HOST=${var.db_hostname}",
    "DOCKER_DB_USER=postgres",
  ]
  networks_advanced {
    name = "${var.db_network}"
  }
  networks_advanced {
    name = "${var.public_network}"
  }
  networks_advanced {
    name = "${docker_network.idp_private.name}"
  }
  user = "appuser"
#  volumes {
#    host_path = "${var.host_path}"
#    container_path = "/upaya"
#  }
}

resource "docker_container" "redis" {
  name = "redis"
  hostname = "redis"
  image = "${docker_image.redis.latest}"
  networks_advanced {
    name = "${docker_network.idp_private.name}"
  }
}

resource "docker_image" "redis" {
  name = "redis"
}

resource "docker_network" "idp_private" {
  name = "idp_private"
  internal = true
}
