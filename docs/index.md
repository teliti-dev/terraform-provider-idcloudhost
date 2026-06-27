---
page_title: "idcloudhost Provider"
description: |-
  Terraform provider untuk IDCloudHost Cloud VPS. Supports VM, network, floating IP, firewall, load balancer, dan object storage.
---

# idcloudhost Provider

Provider untuk mengelola resources di [IDCloudHost](https://idcloudhost.com).

## Konfigurasi

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
  auth_token = var.auth_token  # atau via env: IDCLOUDHOST_AUTH_TOKEN
  region     = "jkt01"         # jkt01 | jkt03 | sgp01
}
```

## Schema

### Optional

- `auth_token` (String, Sensitive) — API token IDCloudHost. Bisa juga via env `IDCLOUDHOST_AUTH_TOKEN`.
- `region` (String) — Region IDCloudHost. Default: `jkt01`. Pilihan: `jkt01`, `jkt03`, `sgp01`.
