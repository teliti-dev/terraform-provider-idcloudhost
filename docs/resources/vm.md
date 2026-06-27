---
page_title: "idcloudhost_vm Resource - terraform-provider-idcloudhost"
description: |-
  Membuat dan mengelola Virtual Machine di IDCloudHost.
---

# idcloudhost_vm (Resource)

Membuat Virtual Machine di IDCloudHost.

> **Catatan:** VM tidak mendapat public IP secara otomatis. Gunakan `idcloudhost_floating_ip` untuk public IP, lalu assign manual di dashboard IDCloudHost.

## Contoh

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
  initial_password   = "P@ssw0rd123!"
  public_key         = file("~/.ssh/id_ed25519.pub")
  billing_account_id = var.billing_account_id
  backup             = false

  # Pilih server class — kosongkan untuk DEFAULT pool region
  designated_pool_uuid = "00000000-0000-0000-0000-000022006840"  # Basic jkt01

  # Attach ke private network saat create
  network_uuid = idcloudhost_network.main.uuid
}
```

## Server Class (designated_pool_uuid)

| Region | Pool | Keterangan | UUID |
|---|---|---|---|
| `jkt01` | Basic | Standard | `00000000-0000-0000-0000-000022006840` |
| `jkt01` | Intel eXtreme | 5x Faster | `9b6bf39f-6559-4e06-be68-6252e980468d` |
| `jkt01` | **AMD eXtreme** | 6x Faster — **DEFAULT** | `6d4026f6-1a7b-4f32-966b-2e739d4533b1` |
| `jkt03` | **Basic** | Standard — **DEFAULT** | `1bcdc355-83b9-41db-83f4-7162b19a2990` |
| `sgp01` | **Intel Pro** | 3x Faster — **DEFAULT** | `e2ab9e01-43ef-4a20-93e2-30a40d7545fb` |

## Schema

### Required

- `billing_account_id` (Number) — ID billing account IDCloudHost.
- `disks` (Number) — Ukuran disk dalam GB. Min: 20, Max: 240.
- `initial_password` (String, Sensitive) — Password awal VM. Min 8 karakter, harus ada huruf besar, angka, dan simbol.
- `memory` (Number) — RAM dalam MB. Min: 2048, Max: 65536.
- `name` (String) — Nama VM.
- `os_name` (String) — Nama OS. Contoh: `ubuntu`.
- `os_version` (String) — Versi OS. Contoh: `20.04`, `22.04`.
- `username` (String) — Username untuk login ke VM.
- `vcpu` (Number) — Jumlah vCPU. Min: 2, Max: 32.

### Optional

- `backup` (Boolean) — Aktifkan auto backup. Default: `false`.
- `description` (String) — Deskripsi VM.
- `designated_pool_uuid` (String) — UUID server class (pool). Kosongkan untuk DEFAULT pool region. Lihat tabel Server Class di atas.
- `network_uuid` (String) — UUID private network untuk di-attach saat VM dibuat.
- `public_key` (String) — SSH public key yang di-inject ke VM via cloud-init.
- `source_replica` (String) — Nama replica untuk clone dari backup.
- `source_uuid` (String) — UUID VM sumber untuk clone.
- `timeouts` (Block, Optional) — Timeout konfigurasi.
  - `create` (String) — Timeout create VM. Default: `5m`.

### Read-Only

- `created_at` (String)
- `hostname` (String)
- `hypervisor_id` (String)
- `id` (String) — ID resource.
- `mac` (String)
- `private_ipv4` (String) — IP private VM di dalam network.
- `status` (String) — Status VM: `running`, `stopped`, dll.
- `storage` (List of Object) — Disk yang ter-attach ke VM.
- `tags` (List of String)
- `updated_at` (String)
- `user_id` (Number)
- `uuid` (String) — UUID unik VM, dipakai untuk referensi ke resource lain.
