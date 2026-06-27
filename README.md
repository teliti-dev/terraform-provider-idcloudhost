# Terraform Provider — IDCloudHost

Custom Terraform provider untuk [IDCloudHost](https://idcloudhost.com) Cloud VPS.

Fork dari [bapung/terraform-provider-idcloudhost](https://github.com/bapung/terraform-provider-idcloudhost), dimodifikasi oleh [teliti-dev](https://github.com/teliti-dev).

## Instalasi (Manual)

Provider ini tidak ada di Terraform Registry — harus di-build dan install manual.

```bash
git clone git@github.com:teliti-dev/tf-idcloudhost.git
cd tf-idcloudhost

OS=$(go env GOOS)
ARCH=$(go env GOARCH)
PLUGIN_DIR=~/.terraform.d/plugins/teliti-dev/idcloudhost/0.2.0/${OS}_${ARCH}

go build -o terraform-provider-idcloudhost .
mkdir -p "$PLUGIN_DIR"
cp terraform-provider-idcloudhost "$PLUGIN_DIR/"
```

## Penggunaan

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
  # auth_token dibaca dari env: IDCLOUDHOST_AUTH_TOKEN
  # region default: jkt01
}
```

## Regions

| Region | Lokasi |
|---|---|
| `jkt01` | Jakarta 1 |
| `jkt03` | Jakarta 3 |
| `sgp01` | Singapore 1 |

## Server Class (Pool UUID)

Diisi ke field `designated_pool_uuid` pada resource `idcloudhost_vm`. Kosongkan untuk pakai DEFAULT.

| Region | Pool | Keterangan | UUID |
|---|---|---|---|
| `jkt01` | Basic | Standard | `00000000-0000-0000-0000-000022006840` |
| `jkt01` | Intel eXtreme | 5x Faster | `9b6bf39f-6559-4e06-be68-6252e980468d` |
| `jkt01` | **AMD eXtreme** | 6x Faster — **DEFAULT** | `6d4026f6-1a7b-4f32-966b-2e739d4533b1` |
| `jkt03` | **Basic** | Standard — **DEFAULT** | `1bcdc355-83b9-41db-83f4-7162b19a2990` |
| `sgp01` | **Intel Pro** | 3x Faster — **DEFAULT** | `e2ab9e01-43ef-4a20-93e2-30a40d7545fb` |

Semua pool limit: vCPU 2–32, RAM 2048–65536 MB, Disk 20–1000 GB.

## Contoh: VM + Floating IP + Private Network

```hcl
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
  initial_password   = "P@ssw0rd123!"
  public_key         = file("~/.ssh/id_ed25519.pub")
  billing_account_id = var.billing_account_id
  backup             = false

  # Pilih server class (opsional — kosongkan untuk DEFAULT)
  designated_pool_uuid = "00000000-0000-0000-0000-000022006840"  # Basic jkt01

  # Attach ke private network saat create
  network_uuid = idcloudhost_network.main.uuid
}

resource "idcloudhost_floating_ip" "web" {
  name               = "web-1-ip"
  billing_account_id = var.billing_account_id
}

output "public_ip"  { value = idcloudhost_floating_ip.web.address }
output "private_ip" { value = idcloudhost_vm.web.private_ipv4 }
```

> **Catatan:** Floating IP harus di-assign ke VM secara manual di dashboard IDCloudHost setelah `terraform apply`.

## Resources

| Resource | Keterangan |
|---|---|
| `idcloudhost_vm` | Virtual machine |
| `idcloudhost_network` | Private network |
| `idcloudhost_floating_ip` | Public IP |
| `idcloudhost_firewall` | Firewall rules |
| `idcloudhost_loadbalancer` | Load balancer |
| `idcloudhost_objectstorage` | Object storage bucket |
| `idcloudhost_vm_disks` | Tambah disk ke VM |

## Data Sources

| Data Source | Keterangan |
|---|---|
| `idcloudhost_vms` | List semua VM yang ada |
