resource "docker_container" "idp" {
  name = "web"
  image = "identity-idp_web:latest"
  ports {
    internal = 3000
    external = 3000
  }
  env = [
    "REDIS_URL=redis://redis",
    "DATABASE_URL=postgres://postgres@db",
    "DOCKER_DB_HOST=db",
    "DOCKER_DB_USER=postgres",
  ]
  networks_advanced {
    name = "${docker_network.private.name}"
  }
  user = "upaya"
  volumes {
    host_path = "/Users/usds/Development/identity-idp"
    container_path = "/upaya"
  }
}

resource "docker_container" "postgres" {
  name = "db"
  image = "${docker_image.postgres.latest}"
  networks_advanced {
    name = "${docker_network.private.name}"
    aliases = ["db"]
  }
  publish_all_ports=true
}

resource "docker_container" "redis" {
  name = "redis"
  image = "${docker_image.redis.latest}"
  networks_advanced {
    name = "${docker_network.private.name}"
    aliases = ["redis"]
  }
}

resource "docker_image" "postgres" {
  name = "postgres"
}

resource "docker_image" "redis" {
  name = "redis"
}

resource "docker_network" "private" {
  name = "private"
}
