# Terraform Provider — IDCloudHost

[![Terraform Registry](https://img.shields.io/badge/Terraform%20Registry-teliti--dev%2Fidcloudhost-623CE4)](https://registry.terraform.io/providers/teliti-dev/idcloudhost/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/teliti-dev/terraform-provider-idcloudhost)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Terraform provider untuk [IDCloudHost](https://idcloudhost.com) Cloud VPS. Kelola virtual machine, private network, floating IP, firewall, load balancer, dan object storage melalui Terraform.

Diturunkan dari [bapung/terraform-provider-idcloudhost](https://github.com/bapung/terraform-provider-idcloudhost) dan dikelola oleh [teliti-dev](https://github.com/teliti-dev).

---

## Persyaratan

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (hanya untuk development lokal)
- Akun [IDCloudHost](https://idcloudhost.com) dengan API token

---

## Instalasi

Tambahkan provider ke konfigurasi Terraform:

```hcl
terraform {
  required_providers {
    idcloudhost = {
      source  = "teliti-dev/idcloudhost"
      version = "~> 0.2"
    }
  }
}
```

Kemudian jalankan:

```bash
terraform init
```

---

## Autentikasi

Set API token melalui environment variable (direkomendasikan):

```bash
export IDCLOUDHOST_AUTH_TOKEN="api-token-anda"
```

Atau langsung di blok provider (tidak direkomendasikan untuk production):

```hcl
provider "idcloudhost" {
  auth_token = "api-token-anda"
  region     = "jkt01"
}
```

---

## Region

| Region | Lokasi      |
|--------|-------------|
| `jkt01` | Jakarta 1  |
| `jkt03` | Jakarta 3  |
| `sgp01` | Singapore 1 |

---

## Server Class

Gunakan field `designated_pool_uuid` pada `idcloudhost_vm` untuk memilih server class. Kosongkan untuk menggunakan DEFAULT region.

| Region | Pool           | Keterangan        | UUID                                   |
|--------|----------------|-------------------|----------------------------------------|
| `jkt01` | Basic         | Standard          | `00000000-0000-0000-0000-000022006840` |
| `jkt01` | Intel eXtreme | 5x Lebih Cepat    | `9b6bf39f-6559-4e06-be68-6252e980468d` |
| `jkt01` | **AMD eXtreme** | 6x Lebih Cepat — **DEFAULT** | `6d4026f6-1a7b-4f32-966b-2e739d4533b1` |
| `jkt03` | **Basic**     | Standard — **DEFAULT** | `1bcdc355-83b9-41db-83f4-7162b19a2990` |
| `sgp01` | **Intel Pro** | 3x Lebih Cepat — **DEFAULT** | `e2ab9e01-43ef-4a20-93e2-30a40d7545fb` |

Semua pool mendukung: vCPU 2–32, RAM 2048–65536 MB, Disk 20–1000 GB.

---

## Contoh Penggunaan

### VM dengan Private Network dan Floating IP

```hcl
terraform {
  required_providers {
    idcloudhost = {
      source  = "teliti-dev/idcloudhost"
      version = "~> 0.2"
    }
  }
}

provider "idcloudhost" {
  region = "jkt01"
}

resource "idcloudhost_network" "main" {
  name = "my-network"
}

resource "idcloudhost_vm" "web" {
  name               = "web-1"
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

  # Opsional: pilih server class (kosongkan untuk DEFAULT region)
  designated_pool_uuid = "00000000-0000-0000-0000-000022006840"

  # Opsional: attach ke private network saat VM dibuat
  network_uuid = idcloudhost_network.main.uuid
}

resource "idcloudhost_floating_ip" "web" {
  name               = "web-1-ip"
  billing_account_id = var.billing_account_id
}

output "public_ip"  { value = idcloudhost_floating_ip.web.address }
output "private_ip" { value = idcloudhost_vm.web.private_ipv4 }
```

> **Catatan:** Floating IP harus di-assign ke VM secara manual melalui dashboard IDCloudHost setelah `terraform apply`. Provider hanya mereservasi IP, tidak melakukan auto-assign.

---

## Resources

| Resource | Keterangan |
|---|---|
| [`idcloudhost_vm`](docs/resources/vm.md) | Virtual machine |
| [`idcloudhost_network`](docs/resources/network.md) | Private network |
| [`idcloudhost_floating_ip`](docs/resources/floating_ip.md) | Public floating IP |
| [`idcloudhost_firewall`](docs/resources/firewall.md) | Firewall rules |
| [`idcloudhost_loadbalancer`](docs/resources/loadbalancer.md) | Load balancer |
| [`idcloudhost_objectstorage`](docs/resources/objectstorage.md) | Object storage bucket |
| [`idcloudhost_vm_disks`](docs/resources/vm_disks.md) | Disk tambahan untuk VM |

## Data Sources

| Data Source | Keterangan |
|---|---|
| [`idcloudhost_vms`](docs/data-sources/vms.md) | List semua VM di akun |

---

## Development Lokal

```bash
git clone git@github.com:teliti-dev/terraform-provider-idcloudhost.git
cd terraform-provider-idcloudhost

# Build dan install lokal
make install

# Jalankan tests
make test
```

Perintah `make install` mem-build binary dan menyalinnya ke `~/.terraform.d/plugins/` agar bisa digunakan dengan konfigurasi Terraform lokal.

---

## Kontribusi

Pull request diterima. Untuk perubahan besar, buka issue terlebih dahulu.

---

## Lisensi

[MIT](LICENSE)
