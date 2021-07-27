# request detailt for a workspace with the name "Some workspace"
data "ardoq_workspace" "someworkspace" {
    name = "Some workspace"
}

# returns all references in workspace
data "ardoq_references" "all" {
  # use the retrieved data from above
  root_workspace = data.ardoq_workspace.someworkspace.id
}

# output all components
output "all_references" {
  value = data.ardoq_references.all.components
}