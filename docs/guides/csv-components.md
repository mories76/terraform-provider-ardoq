---
subcategory: ""
page_title: "Create components from a CSV - Ardoq Provider"
description: |-
    An example how to use a CSV to create components.
---

# Create components from a CSV

Given a CSV file with the following content:

```csv
name,description
Component 1,And the description
Component 2,Another description
```

You could create/manage a `arodq_component` for every row address in the CSV with the following config:

```terraform
locals {
  componentscsv = csvdecode(file("${path.module}/guide.csv"))
  components    = { for component in local.components : component.name => component }
}

resource "ardoq_component" "component" {
  for_each = local.components

  name = each.key
  description = each.value.description
}
```