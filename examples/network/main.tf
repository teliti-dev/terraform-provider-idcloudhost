terraform {
  required_providers {
    idcloudhost = {
      version = "0.2.0"
      source  = "github.com/bapung/idcloudhost"
    }
  }
}

provider "idcloudhost" {
  # auth_token is read from IDCLOUDHOST_AUTH_TOKEN environment variable
  # region defaults to "jkt01"
}

# Create a private network
resource "idcloudhost_network" "example_network" {
  name    = "my-private-network"
  default = false
}

# Create a network and set it as default
resource "idcloudhost_network" "default_network" {
  name    = "my-default-network"
  default = true
}

# Output network details
output "example_network_uuid" {
  value       = idcloudhost_network.example_network.uuid
  description = "The UUID of the example network"
}

output "example_network_id" {
  value       = idcloudhost_network.example_network.id
  description = "The ID of the example network"
}

output "default_network_uuid" {
  value       = idcloudhost_network.default_network.uuid
  description = "The UUID of the default network"
}
