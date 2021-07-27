# request details for a workspace with the name "Some workspace"
data "ardoq_workspace" "someworkspace" {
    name = "Some workspace"
}

# returns all workspaces in workspace
data "ardoq_components" "all" {
  # use the retrieved data from above
  root_workspace = data.ardoq_workspace.someworkspace.id
}

# output all components
output "all_components" {
  value = data.ardoq_components.all.components
}