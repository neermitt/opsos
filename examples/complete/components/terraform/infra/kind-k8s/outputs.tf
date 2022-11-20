output "name" {
  value       = module.kind_cluster.name
  description = "the name of the kind_cluster"
}

output "kubeconfig" {
  value       = module.kind_cluster.kubeconfig
  description = "the kubeconfig of the kind_cluster"
}

output "ipAddress" {
  value       = module.kind_cluster.ipAddress
  description = "the IP address of the kind_cluster controller node"
}

output "kubeconfig_with_ip" {
  value       = module.kind_cluster.kubeconfig_with_ip
  description = "the kubeconfig of the kind_cluster with the IP address of the controller node"
}

output "context_name" {
  value       = module.kind_cluster.context_name
  description = "context name for the kind cluster in kube_config file"
}

output "loadbalancer_container_name" {
  value       = module.kind_cluster.loadbalancer_container_name
  description = "loadbalancer container name for the kind cluster"
}