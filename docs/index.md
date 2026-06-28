---
page_title: "idcloudhost Provider"
description: |-
  Use the IDCloudHost provider to manage cloud infrastructure resources including virtual machines, private networks, floating IPs, firewalls, load balancers, and object storage.
---

# IDCloudHost Provider

The IDCloudHost provider is used to manage resources in [IDCloudHost](https://idcloudhost.com) Cloud VPS.

## Authentication

The provider requires an API token. Set it via environment variable (recommended):

```bash
export IDCLOUDHOST_AUTH_TOKEN="your-api-token"
```

## Example Usage

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
  # auth_token is read from IDCLOUDHOST_AUTH_TOKEN env var
  region = "jkt01"
}
```

## Schema

### Optional

- `auth_token` (String, Sensitive) — IDCloudHost API token. Can also be set via the `IDCLOUDHOST_AUTH_TOKEN` environment variable.
- `region` (String) — IDCloudHost region. Defaults to `jkt01`. Valid values: `jkt01`, `jkt03`, `sgp01`.
