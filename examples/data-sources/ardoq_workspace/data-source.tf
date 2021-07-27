# request details for a workspace with the name "Some workspace"
data "ardoq_workspace" "someworkspace" {
    name = "Some workspace"
}

output "workspace_output" {
  value = data.ardoq_workspace.someworkspace
}