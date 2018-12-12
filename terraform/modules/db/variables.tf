variable "name" {}

variable "db_name" {
  description = "Name of database to initialize on startup"
}

variable "host_path" {
  description = "Host path to use for persisting data"
}

variable "use_local" {
  description = "If set to true, use a locally running database, otherwise, create a new one and private network."
}
