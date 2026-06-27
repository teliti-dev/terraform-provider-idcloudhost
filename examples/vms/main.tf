terraform {
  required_providers {
    idcloudhost = {
      version = "0.2.0"
      source  = "teliti-dev/idcloudhost"
    }
  }
}

data "idcloudhost_vms" "all" {}

output "all_vms" {
  value = data.idcloudhost_vms.all
}