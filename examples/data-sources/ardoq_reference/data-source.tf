# returns a reference
data "ardoq_references" "myref" {
  ## TODO
}

# output all references
output "all_references" {
  value = data.ardoq_references.all.references
}