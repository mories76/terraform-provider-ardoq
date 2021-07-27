locals {
  componentscsv = csvdecode(file("${path.module}/guide.csv"))
  components    = { for component in local.components : component.name => component }
}

resource "ardoq_component" "component" {
  for_each = local.components

  name = each.key
  description = each.value.description
}