---
page_title: "jamfpro_jamf_protect"
description: |-
  
---

# jamfpro_jamf_protect (Resource)


## Example Usage
```terraform
resource "jamfpro_jamf_protect" "settings" {
  protect_url  = "https://myinstance.protect.jamfcloud.com/graphql"
  client_id    = "supersecretclientid"
  password     = "supersecretpassword"
  auto_install = true

  timeouts {
    create = "90s"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `client_id` (String, Sensitive) The API client ID for Jamf Protect authentication
- `password` (String, Sensitive) The password for Jamf Protect authentication
- `protect_url` (String, Sensitive) The URL of the Jamf Protect instance

### Optional

- `auto_install` (Boolean) Whether to automatically install Jamf Protect on devices
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `api_client_name` (String) Name of the API client used for integration
- `id` (String) The ID of this resource.
- `last_sync_time` (String) Timestamp of the last successful sync
- `registration_id` (String) Registration ID of the Jamf Protect integration
- `sync_status` (String) Current sync status of the Jamf Protect integration

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `read` (String)
- `update` (String)