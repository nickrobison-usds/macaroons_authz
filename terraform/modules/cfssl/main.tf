resource "docker_container" "cfssl" {
  name = "cfssl"
  image = "nickrobison.com/cfssl:latest"
  hostname = "${var.name}"
  ports {
    internal = 8888
    external = 8888
  }
  networks_advanced {
    name = "${var.public_network}"
  }
}
