terraform {
  required_version = ">= 1.0.0"
  experiments      = [module_variable_optional_attrs]

  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = ">= 0.0.14"
    }
    docker-utils = {
      source  = "Kaginari/docker-utils"
      version = ">= 0.0.5"
    }
    utils = {
      source  = "cloudposse/utils"
      version = ">= 1.5.0"
    }
  }
}