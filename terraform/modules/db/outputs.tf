output "hostname" {
  value = "${local.hostname}"
}

output "network" {
  value = "${var.use_local == 1 ? local.local_network : local.network_name}"
}
