---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ardoq_reference Resource - terraform-provider-ardoq"
subcategory: ""
description: |-
  ardoq_reference resource lets you create a reference
---

# ardoq_reference (Resource)

`ardoq_reference` resource lets you create a reference



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **root_workspace** (String) Id of the source component's workspace
- **source** (String) Id of the source component
- **target** (String) Id of the target component
- **target_workspace** (String) Id of the target component's workspace
- **type** (Number) Type (as defined by the model) i.e. Synchronous, Implicit etc.

### Optional

- **description** (String) Text field describing the reference
- **display_text** (String) Short label describing the reference, is visible in some visualizations
- **fields** (Map of String) All custom fields from the model end up here

### Read-Only

- **id** (String) The unique ID of the reference


