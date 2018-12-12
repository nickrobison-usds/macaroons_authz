resource "docker_container" "cfssl" {
  name = "cfssl"
  image = "nickrobison.com/cfssl:latest"
  hostname = "${var.name}"
  networks_advanced {
    name = "${var.public_network}"
  }
}
