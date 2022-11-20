

module "kind_cluster" {
  source = "git::github.com/neermitt/terraform-kind-k8s.git?ref=0.1.0"

  nodes      = var.nodes
  networking = var.networking
  context    = module.this.context
}
