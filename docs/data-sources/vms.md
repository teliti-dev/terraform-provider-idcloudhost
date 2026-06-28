---
page_title: "idcloudhost_vms Data Source - terraform-provider-idcloudhost"
description: |-
  List all virtual machines in the IDCloudHost account.
---

# idcloudhost_vms (Data Source)

Lists all virtual machines in the IDCloudHost account.

## Example Usage

```hcl
data "idcloudhost_vms" "all" {}

output "vm_list" {
  value = data.idcloudhost_vms.all.vms[*].name
}
```

## Schema

### Read-Only

- `id` (String) — Resource ID.
- `vms` (List of Object) — List of virtual machines. See [Nested Schema](#nested-schema-for-vms) below.

### Nested Schema for `vms`

Read-Only:

- `backup` (Boolean) — Whether automatic backups are enabled.
- `billing_account_id` (Number) — Billing account ID.
- `created_at` (String) — VM creation timestamp.
- `description` (String) — VM description.
- `hostname` (String) — VM hostname.
- `hypervisor_id` (String) — Hypervisor host identifier.
- `id` (Number) — VM numeric ID.
- `mac` (String) — MAC address.
- `memory` (Number) — RAM in MB.
- `name` (String) — VM name.
- `os_name` (String) — Operating system name.
- `os_version` (String) — Operating system version.
- `private_ipv4` (String) — Private IP address.
- `status` (String) — VM status: `running`, `stopped`, etc.
- `storage` (List of Object) — Attached disks. See [Nested Schema](#nested-schema-for-vmsstorage) below.
- `tags` (List of String) — VM tags.
- `updated_at` (String) — Last update timestamp.
- `user_id` (Number) — IDCloudHost user ID.
- `username` (String) — Login username.
- `uuid` (String) — Unique VM UUID.
- `vcpu` (Number) — Number of vCPUs.

### Nested Schema for `vms.storage`

Read-Only:

- `created_at` (String) — Disk creation timestamp.
- `id` (Number) — Disk numeric ID.
- `name` (String) — Disk name.
- `pool` (String) — Storage pool.
- `primary` (Boolean) — Whether this is the primary (root) disk.
- `replica` (List of String) — Replica names (for backed-up disks).
- `shared` (Boolean) — Whether the disk is shared.
- `size` (Number) — Disk size in GB.
- `type` (String) — Disk type.
- `updated_at` (String) — Last update timestamp.
- `user_id` (Number) — IDCloudHost user ID.
- `uuid` (String) — Unique disk UUID.
