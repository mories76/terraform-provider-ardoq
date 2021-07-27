# returns all workspaces in workspace
data "ardoq_models" "all" {
}

# output all components
output "all_models" {
  value = data.ardoq_models.all
}