---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "leaseweb Provider"
subcategory: ""
description: |-
  
---

# leaseweb Provider



## Example Usage

```terraform
# Configuration-based authentication
provider "leaseweb" {
  host  = "127.0.0.1:4010"
  token = "super-secret-token-value"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `token` (String, Sensitive)

### Optional

- `host` (String)
- `scheme` (String)
