############
# Required #
############

variable image_version {
  type = string
  description = "The label of the image to run."
}

variable deployment_args {
  type = list(string)
  default = []
  description = "Command line arguments to pass to residential-dns."
}

############
# Optional #
############

variable create_namespace {
  type = bool
  default = true
  description = "If true, a new namespace will be created with the name set to the value of the namespace_name variable. If false, it will look up an existing namespace with the name of the value of the namespace_name variable."
}

variable namespace_name {
  type = string
  default = "residential-dns"
  description = "The name of the namespace to create or look up."
}

variable env_vars {
  type = list(object({
    name = string
    value = string
  }))
  default = []
  description = "Environment variables to run residential-dns with."
}

variable deployment_annotations {
  type = map
  default = {}
  description = "The set of annotations to add to the deployment pods."
}

variable service_account_annotations {
  type = map
  default = {}
  description = "The set of annotations to add to the service account."
}