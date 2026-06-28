---
page_title: "idcloudhost_floating_ip Resource - terraform-provider-idcloudhost"
description: |-
  Reserve a public floating IP address on IDCloudHost.
---

# idcloudhost_floating_ip (Resource)

Reserves a public floating IP address on IDCloudHost.

-> **Note:** This resource only **reserves** the IP. Assigning it to a VM must be done manually via the IDCloudHost dashboard after `terraform apply`.

## Example Usage

```hcl
resource "idcloudhost_floating_ip" "example" {
  name               = "my-vm-ip"
  billing_account_id = var.billing_account_id
}

output "public_ip" {
  value = idcloudhost_floating_ip.example.address
}
```

## Schema

### Required

- `billing_account_id` (Number) — IDCloudHost billing account ID.
- `name` (String) — Floating IP name.

### Optional

- `assigned_to` (String) — UUID of the VM this IP is assigned to.
- `timeouts` (Block, Optional) — Timeout configuration.
  - `create` (String) — Floating IP creation timeout.

### Read-Only

- `address` (String) — The reserved public IP address.
- `created_at` (String) — Creation timestamp.
- `enabled` (Boolean) — Whether the floating IP is enabled.
- `id` (String) — Resource ID.
- `network_id` (String) — Network ID this IP belongs to.
- `type` (String) — IP type.
- `updated_at` (String) — Last update timestamp.
- `user_id` (Number) — IDCloudHost user ID.
