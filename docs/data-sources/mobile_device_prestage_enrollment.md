---
page_title: "jamfpro_mobile_device_prestage_enrollment"
description: |-
  
---

# jamfpro_mobile_device_prestage_enrollment (Data Source)


## Example Usage
```terraform
data "jamfpro_mobile_device_prestage_enrollment" "example" {
  id = "1"
}

output "prestage_name" {
  value = data.jamfpro_mobile_device_prestage_enrollment.example.display_name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `display_name` (String) The display name of the mobile device prestage.
- `id` (String) The unique identifier of the mobile device prestage.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String)