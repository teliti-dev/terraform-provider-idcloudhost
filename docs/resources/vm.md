---
page_title: "idcloudhost_vm Resource - terraform-provider-idcloudhost"
description: |-
  Create and manage a Virtual Machine on IDCloudHost.
---

# idcloudhost_vm (Resource)

Creates and manages a Virtual Machine on IDCloudHost.

-> **Note:** VMs are created without a public IP by default. Use [`idcloudhost_floating_ip`](floating_ip.md) to reserve a public IP, then assign it to the VM manually via the IDCloudHost dashboard.

## Example Usage

### OS (Linux Distribution)

```hcl
resource "idcloudhost_network" "main" {
  name = "my-network"
}

resource "idcloudhost_vm" "example" {
  name               = "my-vm"
  os_name            = "ubuntu"
  os_version         = "22.04"
  vcpu               = 2
  memory             = 2048
  disks              = 20
  username           = "ubuntu"
  initial_password   = var.vm_password
  public_key         = file("~/.ssh/id_ed25519.pub")
  billing_account_id = var.billing_account_id
  backup             = false

  designated_pool_uuid = "00000000-0000-0000-0000-000022006840"
  network_uuid         = idcloudhost_network.main.uuid
}
```

### App Catalog (Pre-installed Application)

```hcl
resource "idcloudhost_vm" "wordpress" {
  name               = "my-wordpress"
  os_name            = "wordpress"
  os_version         = "6.2"
  vcpu               = 2
  memory             = 2048
  disks              = 20
  username           = "ubuntu"
  initial_password   = var.vm_password
  billing_account_id = var.billing_account_id
}
```

### Empty Disk (No OS)

```hcl
resource "idcloudhost_vm" "blank" {
  name               = "my-blank-vm"
  vcpu               = 2
  memory             = 2048
  disks              = 20
  username           = "admin"
  initial_password   = var.vm_password
  billing_account_id = var.billing_account_id
}
```

## Server Classes (`designated_pool_uuid`)

Leave `designated_pool_uuid` empty to use the **DEFAULT** pool for the configured region.

| Region | Pool           | Description             | UUID                                   |
|--------|----------------|-------------------------|----------------------------------------|
| `jkt01` | Basic         | Standard                | `00000000-0000-0000-0000-000022006840` |
| `jkt01` | Intel eXtreme | 5x Faster               | `9b6bf39f-6559-4e06-be68-6252e980468d` |
| `jkt01` | **AMD eXtreme** | 6x Faster — **DEFAULT** | `6d4026f6-1a7b-4f32-966b-2e739d4533b1` |
| `jkt03` | **Basic**     | Standard — **DEFAULT**  | `1bcdc355-83b9-41db-83f4-7162b19a2990` |
| `sgp01` | **Intel Pro** | 3x Faster — **DEFAULT** | `e2ab9e01-43ef-4a20-93e2-30a40d7545fb` |

All pools support: vCPU 2–32, RAM 2048–65536 MB, Disk 20–1000 GB.

## Schema

### Required

- `billing_account_id` (Number) — IDCloudHost billing account ID.
- `disks` (Number) — Root disk size in GB. Minimum: `20`, Maximum: `240`.
- `initial_password` (String, Sensitive) — Initial VM password. Must be at least 8 characters and contain uppercase letters, numbers, and symbols.
- `memory` (Number) — RAM in MB. Minimum: `2048`, Maximum: `65536`.
- `name` (String) — VM name.
- `username` (String) — Login username for the VM.
- `vcpu` (Number) — Number of vCPUs. Minimum: `2`, Maximum: `32`.

### Optional

- `backup` (Boolean) — Enable automatic backups. Defaults to `false`.
- `description` (String) — VM description.
- `designated_pool_uuid` (String) — Server class (pool) UUID. Leave empty to use the region default. See the [Server Classes](#server-classes-designated_pool_uuid) table above.
- `network_uuid` (String) — UUID of the private network to attach the VM to on creation.
- `os_name` (String) — OS or App Catalog name. Examples: `ubuntu`, `debian`, `centos`, `wordpress`, `lamp`. Leave empty for an empty disk.
- `os_version` (String) — OS or App Catalog version. Examples: `22.04`, `11`, `6.2`. Leave empty for an empty disk.
- `public_key` (String) — SSH public key injected into the VM via cloud-init.
- `source_replica` (String) — Replica name to clone from a backup.
- `source_uuid` (String) — UUID of the source VM to clone.
- `timeouts` (Block, Optional) — Timeout configuration.
  - `create` (String) — VM creation timeout. Defaults to `5m`.

### Read-Only

- `created_at` (String) — VM creation timestamp.
- `hostname` (String) — VM hostname.
- `hypervisor_id` (String) — Hypervisor host identifier.
- `id` (String) — Resource ID.
- `mac` (String) — MAC address.
- `private_ipv4` (String) — Private IP address within the attached network.
- `status` (String) — VM status: `running`, `stopped`, etc.
- `storage` (List of Object) — Disks attached to the VM.
- `tags` (List of String) — VM tags.
- `updated_at` (String) — Last update timestamp.
- `user_id` (Number) — IDCloudHost user ID.
- `uuid` (String) — Unique VM UUID. Use this to reference the VM in other resources.
