resource "docker_container" "idp" {
  name = "web"
  image = "identity-idp_web:latest"
  hostname = "idp_web"
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
  hostname = "db"
  image = "${docker_image.postgres.latest}"
  networks_advanced {
    name = "${docker_network.private.name}"
  }
  volumes {
    host_path = "/Users/usds/Development/identity-idp/db_data"
    container_path = "/var/lib/postgresql/data"
  }
}

resource "docker_container" "redis" {
  name = "redis"
  hostname = "redis"
  image = "${docker_image.redis.latest}"
  networks_advanced {
    name = "${docker_network.private.name}"
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
