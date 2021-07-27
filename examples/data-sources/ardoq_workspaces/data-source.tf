# returns all workspaces
data "ardoq_workspaces" "all" {
}

# output all workspaces
output "all_workspaces" {
  value = data.ardoq_workspaces.all.workspaces
}