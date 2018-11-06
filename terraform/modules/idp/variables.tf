variable "db_hostname" {}

variable "db_network" {}

variable "public_network" {}

variable "name" {
  description = "Hostname to use for the identity provider"
}

variable "host_path" {
  description = "Location of identity provider web resources"
}
