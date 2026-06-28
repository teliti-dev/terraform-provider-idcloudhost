---
page_title: "idcloudhost_network Resource - terraform-provider-idcloudhost"
description: |-
  Create and manage a private network on IDCloudHost.
---

# idcloudhost_network (Resource)

Creates and manages a private network on IDCloudHost. Use the `uuid` output to attach VMs to the network at creation time via [`idcloudhost_vm.network_uuid`](vm.md).

## Example Usage

```hcl
resource "idcloudhost_network" "main" {
  name = "my-network"
}

resource "idcloudhost_vm" "example" {
  name         = "my-vm"
  network_uuid = idcloudhost_network.main.uuid
  # ... other fields
}
```

## Schema

### Required

- `name` (String) — Network name.

### Optional

- `default` (Boolean) — Set this network as the default network.

### Read-Only

- `created_at` (String) — Network creation timestamp.
- `description` (String) — Network description.
- `id` (String) — Resource ID.
- `updated_at` (String) — Last update timestamp.
- `user_id` (Number) — IDCloudHost user ID.
- `uuid` (String) — Unique network UUID. Use this in `idcloudhost_vm.network_uuid`.
